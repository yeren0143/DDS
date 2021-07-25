package writer

import (
	"log"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	"github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
)

var _ IRTPSWriter = (*StatelessWriter)(nil)

// Class StatelessWriter, specialization of RTPSWriter that manages writers that
// don't keep state of the matched readers.
type StatelessWriter struct {
	rtpsWriterBase
	flowControllers                []flowcontrol.IFlowController
	fixedLocators                  common.LocatorList
	matchedReaders                 []ReaderLocator
	lastIntraprocessSequenceNumber uint64
	unsentChanges                  []ChangeForReader
}

func (statelessWriter *StatelessWriter) UnsentChangeAddedToHistory(aChange *common.CacheChangeT,
	maxBlockingTime common.Time) {
	statelessWriter.Mutex.Lock()
	defer statelessWriter.Mutex.Unlock()

	// Request to send all changes to all readers
	if statelessWriter.hist.GetHistorySize() > 0 {
		// Mark all changes as pending
		// statelessWriter.unsentChanges = make([]ChangeForReader, statelessWriter.hist.GetHistorySize())
		// copy(statelessWriter.unsentChanges, statelessWriter.hist.Changes...)
		log.Fatalln("not impl")
	}
	log.Fatalln("not impl")
}

func (statelessWriter *StatelessWriter) AddFlowController(controller flowcontrol.IFlowController) {
	statelessWriter.flowControllers = append(statelessWriter.flowControllers, controller)
}

func (statelessWriter *StatelessWriter) UnsentChangesReset() {
	statelessWriter.Mutex.Lock()
	defer statelessWriter.Mutex.Unlock()

	// Request to send all changes to all readers
	log.Fatalln("not Impl")

}

func (statelessWriter *StatelessWriter) SetFixedLocators(locators common.LocatorList) bool {
	statelessWriter.Mutex.Lock()
	defer statelessWriter.Mutex.Unlock()
	statelessWriter.fixedLocators.Append(&statelessWriter.fixedLocators)
	for _, loc := range statelessWriter.fixedLocators.Locators {
		statelessWriter.RTPSParticipant.CreateSenderResources(&loc)
	}
	return true
}

func (statelessWriter *StatelessWriter) getBuiltinGuid() {

}

func (statelessWriter *StatelessWriter) SendAnyUnsentChanges() {

}

func (statelessWriter *StatelessWriter) init(parent endpoint.IEndpointParent, att *attributes.WriterAttributes) {
	statelessWriter.getBuiltinGuid()

	locAlloc := parent.GetAttributes().Allocation.Locators
	for i := uint32(0); i < att.MatchedReadersAllocation.Initial; i++ {
		loc := NewReaderLocator(statelessWriter, locAlloc.MaxUnicastLocators, locAlloc.MaxMulticastLocators)
		statelessWriter.matchedReaders = append(statelessWriter.matchedReaders, *loc)
	}
}

func NewStatelessWriter(parent endpoint.IEndpointParent, guid *common.GUIDT,
	att *attributes.WriterAttributes, payloadPool history.IPayloadPool, changePool history.IChangePool,
	hist *history.WriterHistory, wlisten IWriterListener) *StatelessWriter {
	var awriter StatelessWriter
	awriter.rtpsWriterBase = *newRtpsWriterBase(parent, guid, att, payloadPool, changePool, hist, wlisten)
	//changeSize := history.ResourceLimitsFromHistory(hist.Att)
	awriter.hist.Writer = &awriter
	awriter.hist.Mutex = awriter.Mutex
	awriter.init(parent, att)
	return &awriter
}
