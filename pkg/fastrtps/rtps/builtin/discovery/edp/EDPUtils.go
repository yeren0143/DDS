package edp

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
)

type WriterHistoryPair struct {
	AWriter *writer.StatefulWriter
	Hist    *history.WriterHistory
}

type ReaderHistoryPair struct {
	AReader *reader.StatefulReader
	Hist    *history.ReaderHistory
}

func createPayloadPool(topicName string, histAtt *attributes.HistoryAttributes, isReader bool) history.ITopicPayloadPool {
	poolCfg := history.FromHistoryAttributes(histAtt)
	poolProxy := history.GetTopicPayloadPoolProxy(topicName, poolCfg)
	poolProxy.ReserveHistory(poolCfg, isReader)
	return poolProxy
}

func releasePayloadPool(pool history.ITopicPayloadPool, histAtt *attributes.HistoryAttributes, isReader bool) {
	if pool != nil {
		poolConfig := history.FromHistoryAttributes(histAtt)
		pool.ReleaseHistory(poolConfig, isReader)
		history.ReleaseTopicPayloadPool(pool)
	}
}

func createEDPReader(participant protocol.IParticipant, topicName string, entityID *common.EntityIDT,
	histAtt *attributes.HistoryAttributes, ratt *attributes.ReaderAttributes, listener reader.IReaderListener,
	payloadPool history.ITopicPayloadPool, edpReader *ReaderHistoryPair) bool {

	payloadPool = createPayloadPool(topicName, histAtt, true)
	edpReader.Hist = history.NewReaderHistory(histAtt)

	ok, aReader := participant.CreateReader(ratt, payloadPool,
		edpReader.Hist, listener, entityID, true, true)

	if ok {
		edpReader.AReader = aReader.(*reader.StatefulReader)
	}
	if !ok || edpReader.AReader == nil {
		edpReader.Hist = nil
		releasePayloadPool(payloadPool, histAtt, true)
	}

	return ok
}

func createEDPWriter(participant protocol.IParticipant, topicName string, entityID *common.EntityIDT,
	histAtt *attributes.HistoryAttributes, watt *attributes.WriterAttributes, listener writer.IWriterListener,
	payloadPool history.ITopicPayloadPool, edpWriter *WriterHistoryPair) bool {

	payloadPool = createPayloadPool(topicName, histAtt, false)
	edpWriter.Hist = history.NewWriterHistory(histAtt)
	ok, awriter := participant.CreateWriter(watt, payloadPool, edpWriter.Hist, listener, entityID, true)

	if ok {
		edpWriter.AWriter = awriter.(*writer.StatefulWriter)
	} else {
		edpWriter.Hist = nil
		releasePayloadPool(payloadPool, histAtt, false)
	}

	return ok
}
