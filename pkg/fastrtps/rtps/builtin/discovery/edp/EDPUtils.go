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

func createEDPReader(participant *protocol.IParticipant, topicName string, entityID *common.EntityIDT,
	histAtt *attributes.HistoryAttributes, ratt *attributes.ReaderAttributes, listener reader.IReaderListener,
	payloadPool history.ITopicPayloadPool) bool {
	return false
}
