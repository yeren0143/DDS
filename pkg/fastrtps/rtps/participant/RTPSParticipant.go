package participant

import (
	"log"
	"math"
	"os"
	"reflect"
	"sync"

	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/message"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/protocol"
	"github.com/yeren0143/DDS/fastrtps/rtps/endpoint"
	"github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	"github.com/yeren0143/DDS/fastrtps/rtps/history"
	"github.com/yeren0143/DDS/fastrtps/rtps/network"
	"github.com/yeren0143/DDS/fastrtps/rtps/persistence"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
	"github.com/yeren0143/DDS/fastrtps/rtps/transport"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
	"github.com/yeren0143/DDS/fastrtps/utils"
)

type typeCheckFn func(string) bool

// ReceiverControlBlock is a struct we use to encapsulate the resources that take part
// in message reception.  It contains:
// -A ReceiverResource (as produced by the NetworkFactory Element)
// -Its associated MessageReceiver
type ReceiverControlBlock struct {
	ReceiverRes *network.ReceiverResource
	MsgReceiver *message.Receiver
}

//var _ builtin.IProtocolParent = (*RTPSParticipant)(nil)
var _ protocol.IParticipant = (*RTPSParticipant)(nil)
var _ endpoint.IEndpointParent = (*RTPSParticipant)(nil)

//RTPSParticipant allows the creation and removal of writers and readers. It manages the send and receive threads.
type RTPSParticipant struct {
	DomaninID                 uint32
	Att                       *attributes.RTPSParticipantAttributes
	GUID                      common.GUIDT
	PersistenceGUID           common.GUIDT // Persistence guid of the RTPSParticipant
	eventThr                  *resources.ResourceEvent
	BuiltinProtocols          *builtin.Protocols
	ResSemaphore              *utils.Semaphore // Semaphore to wait for the listen thread creation.
	IDCounter                 uint32           // Id counter to correctly assign the ids to writers and readers.
	AllWriterList             []writer.IRTPSWriter
	AllReaderList             []reader.IRTPSReader
	UserWriterList            []writer.IRTPSWriter
	UserReaderList            []reader.IRTPSReader
	Controllers               []flowcontrol.IFlowController
	networkFactory            *network.NetFactory
	asyncThread               *writer.AsyncWriterThread
	Listener                  *RTPSParticipantListener
	IntraProcessOnly          bool
	hasShmTransport           bool
	checkFn                   typeCheckFn
	receiverResourceList      map[*ReceiverControlBlock]bool
	receiverResourceListMutex sync.Mutex
	sendBuffers               *message.SendBuffersManager
	sendResourceList          transport.SenderResourceList
	sendResourceMutex         sync.Mutex
	mutex                     sync.Mutex
}

// TODO:
func shouldBeIntraProcessOnly(pparam *attributes.RTPSParticipantAttributes) bool {
	return false
}

//DidMutationTookPlaceOnMeta Compare metatraffic locators list searching for mutations
func (participant *RTPSParticipant) DidMutationTookPlaceOnMeta(multicast, unicast *common.LocatorList) bool {
	if participant.Att.Builtin.MetatrafficMulticastLocatorList == multicast &&
		participant.Att.Builtin.MetatrafficUnicastLocatorList == unicast {
		return false
	}

	//TODO
	return true
}

func (participant *RTPSParticipant) NetworkFactory() *network.NetFactory {
	return participant.networkFactory
}

func (participant *RTPSParticipant) GetGuid() *common.GUIDT {
	return &participant.GUID
}

//NetworkFactoryHasRegisteredTransports judgment registered transport number
func (participant *RTPSParticipant) NetworkFactoryHasRegisteredTransports() bool {
	return participant.networkFactory.NumberOfRegisteredTransports() > 0
}

func (participant *RTPSParticipant) AssertRemoteParticipantLiveliness(remoteGUID *common.GUIDPrefixT) {
	if participant.BuiltinProtocols != nil && participant.BuiltinProtocols.PDP != nil {
		// TODO:
		// participant.BuiltinProtocols.PDP.AssertRemoteParticipantLiveliness()
	}
}

//Enable Create receiver resources and start builtin protocols
func (participant *RTPSParticipant) Enable() {
	if !participant.BuiltinProtocols.InitBuiltinProtocol(participant, participant.Att.Builtin) {
		log.Fatal("The builtin protocols were not correctly initialized")
	}

	for receiver := range participant.receiverResourceList {
		receiver.ReceiverRes.RegisterReceiver(receiver.MsgReceiver)
	}
}

func (participant *RTPSParticipant) IsIntraprocessOnly() bool {
	return participant.IntraProcessOnly
}

func (participant *RTPSParticipant) SendSync(msg *common.CDRMessage, locators []common.Locator, maxBlockingTimePoint common.Time) bool {
	//retCode := false
	// TODO: try_lock_until
	participant.sendResourceMutex.Lock()
	defer participant.sendResourceMutex.Unlock()
	for _, sendRes := range participant.sendResourceList {
		var locatorList common.LocatorList
		locatorList.Locators = locators
		sendRes.Send(msg.Buffer, locatorList, maxBlockingTimePoint)
	}
	return true

}

func (participant *RTPSParticipant) applyLocatorAdaptRule(loc *common.Locator) *common.Locator {
	// This is a completely made up rule
	// It is transport responsability to interpret this new port.
	loc.Port += uint32(participant.Att.Port.ParticipantIDGain)
	return loc
}

/** Create the new ReceiverResources needed for a new Locator,
    contains the calls to assignEndpointListenResources
	and consequently assignEndpoint2LocatorList
	@param pend - Pointer to the endpoint which triggered the creation of the Receivers
*/
func (participant *RTPSParticipant) createAndAssociateReceiverswithEndpoint(pend endpoint.IEndpoint) bool {
	att := pend.GetAttributes()
	if att.UnicastLocatorList.Length() == 0 && att.MulticastLocatorList.Length() == 0 {
		att.UnicastLocatorList.Append(participant.Att.DefaultUnicastLocatorList)
	}
	participant.createReceiverResources(&att.UnicastLocatorList, false, true)
	participant.createReceiverResources(&att.MulticastLocatorList, false, true)

	// Associate the Endpoint with ReceiverControlBlock
	participant.assignEndpointListenResources(pend)

	return true
}

// Assign an endpoint to the ReceiverResources, based on its LocatorLists.
// RECEIVER RESOURCE METHODS
func (participant *RTPSParticipant) assignEndpointListenResources(endp endpoint.IEndpoint) bool {
	/* No need to check for emptiness on the lists, as it was already done on part function
	   In case are using the default list of Locators they have already been embedded to the parameters
	*/
	participant.assignEndpoint2LocatorList(endp, &endp.GetAttributes().UnicastLocatorList)
	participant.assignEndpoint2LocatorList(endp, &endp.GetAttributes().MulticastLocatorList)
	return true
}

// Assign an endpoint to the ReceiverResources as specified specifically on parameter list
func (participant *RTPSParticipant) assignEndpoint2LocatorList(endp endpoint.IEndpoint, list *common.LocatorList) bool {
	for _, item := range list.Locators {
		participant.receiverResourceListMutex.Lock()
		for rcvControlBlock := range participant.receiverResourceList {
			if rcvControlBlock.ReceiverRes.SupportsLocator(&item) {
				rcvControlBlock.MsgReceiver.AssociateEndpoint(endp)
			}
		}
		participant.receiverResourceListMutex.Unlock()
	}

	return true
}

func (participant *RTPSParticipant) createReceiverResources(locatorList *common.LocatorList,
	applyMutation bool, registerReceiver bool) {
	var newItemsBuffer []*network.ReceiverResource
	maxReceiverBufferSize := ^uint32(0)

	for i, loc := range locatorList.Locators {
		ret, newItems := participant.networkFactory.BuildReceiverResources(&loc, maxReceiverBufferSize)
		if !ret && applyMutation {
			for tries := uint32(0); !ret && tries < participant.Att.Builtin.MutationTries; tries++ {
				locatorList.Locators[i] = *participant.applyLocatorAdaptRule(&loc)
				ret, newItems = participant.networkFactory.BuildReceiverResources(&locatorList.Locators[i], maxReceiverBufferSize)
				newItemsBuffer = append(newItemsBuffer, newItems...)
			}
		} else {
			newItemsBuffer = append(newItemsBuffer, newItems...)
		}

		for _, buffer := range newItemsBuffer {
			participant.receiverResourceListMutex.Lock()
			defer participant.receiverResourceListMutex.Unlock()
			controlBlock := ReceiverControlBlock{
				ReceiverRes: buffer,
				MsgReceiver: nil,
			}
			controlBlock.MsgReceiver = message.NewMessageReceiver(participant, buffer.MaxMessageSize)
			participant.receiverResourceList[&controlBlock] = true
			// start reception
			if registerReceiver {
				controlBlock.ReceiverRes.RegisterReceiver(controlBlock.MsgReceiver)
			}
		}
	}
}

func (participant *RTPSParticipant) existsEntityID(ent *common.EntityIDT, kind common.EndpointKindT) bool {
	if kind == common.KWriter {
		for _, userWriter := range participant.UserWriterList {
			if *ent == userWriter.GetGUID().EntityID {
				return true
			}
		}
	} else {
		for _, userReader := range participant.UserReaderList {
			if *ent == userReader.GetGUID().EntityID {
				return true
			}
		}
	}

	return false
}

func (participant *RTPSParticipant) GetMaxMessageSize() uint32 {
	maxReceiverBufferSize := ^uint32(0)
	maxMsgSizeBetweenTransport := participant.networkFactory.GetMaxMessageSizeBetweenTransports()
	if maxMsgSizeBetweenTransport < maxReceiverBufferSize {
		maxReceiverBufferSize = maxMsgSizeBetweenTransport
	}

	return maxReceiverBufferSize
}

func (participant *RTPSParticipant) GetAttributes() *attributes.RTPSParticipantAttributes {
	return participant.Att
}

func (participant *RTPSParticipant) preprocessEndpointAttributes(debugLable string, entityID *common.EntityIDT,
	att *attributes.EndpointAttributes, kind common.EndpointKindT,
	noKey, withKey common.Octet) (bool, *common.EntityIDT) {
	var entID common.EntityIDT
	if !att.UnicastLocatorList.Valid() {
		log.Fatalln("Unicast Locator List for ", debugLable, " contains invalid Locator")
		return false, nil
	}
	if !att.MulticastLocatorList.Valid() {
		log.Fatalln("Multicast Locator List for ", debugLable, " contains invalid Locator")
		return false, nil
	}
	if !att.RemoteLocatorList.Valid() {
		log.Fatalln("Remote Locator List for ", debugLable, " contains invalid Locator")
		return false, nil
	}

	if entityID == common.KEIDUnknown {
		switch att.TopicKind {
		case common.KNoKey:
			entityID.Value[3] = noKey
		case common.KWithKey:
			entityID.Value[3] = withKey
		default:
			log.Panic("att.TopicKind error")
		}

		var idnum uint32
		if att.GetEntityID() > 0 {
			idnum = uint32(att.GetEntityID())
		} else {
			participant.IDCounter++
			idnum = participant.IDCounter
		}

		entID.Value[2] = common.Octet(idnum)
		entID.Value[1] = common.Octet(idnum >> 8)
		entID.Value[0] = common.Octet(idnum >> 16)
		if participant.existsEntityID(entityID, kind) {
			log.Fatalln("A ", debugLable, " with the same entityId already exists in this RTPSParticipant")
			return false, nil
		}
	} else {
		entID = *entityID
	}

	if att.PersistenceGUID == common.KGuidUnknown {
		// Try to load persistence_guid from property
		persistenceGUIDProperty := attributes.FindProperty(&att.Properties, "dds.persistence.guid")
		if persistenceGUIDProperty != "" {
			// Load persistence_guid from property
			// TODO
			log.Panic("Cannot configure ")
			return false, nil
		}
	}

	return true, &entID
}

func (participant *RTPSParticipant) getPersistenceService(debugLabel string, isBuiltin bool,
	param *attributes.EndpointAttributes) (bool, persistence.IPersistenceService) {
	return true, nil
}

func (participant *RTPSParticipant) CreateSenderResources(locator *common.Locator) {
	participant.sendResourceMutex.Lock()
	defer participant.sendResourceMutex.Unlock()
	participant.networkFactory.BuildSendResources(participant.sendResourceList, locator)
}

func (participant *RTPSParticipant) Wlp() endpoint.IWlp {
	return participant.BuiltinProtocols.WLP
}

func trustedWriter(readerEnt *common.EntityIDT) common.EntityIDT {
	switch readerEnt {
	case common.KEntityIDSPDPReader:
		return *common.KEntityIDSPDPWriter
	case common.KEntityIDSEDPPubReader:
		return *common.KEntityIDSEDPPubWriter
	case common.KEntityIDSEDPSubReader:
		return *common.KEntityIDSEDPSubWriter
	case common.KEntityIDReaderLiveliness:
		return *common.KEntityIDWriterLiveliness
	default:
		return *common.KEIDUnknown
	}
}

func (participant *RTPSParticipant) normalizeEndpointLocators(att *attributes.EndpointAttributes) {
	// Locators with port 0, calculate port.
	for _, loc := range att.UnicastLocatorList.Locators {
		participant.networkFactory.FillDefaultLocatorPort(participant.DomaninID, &loc, participant.Att, false)
	}

	for _, loc := range att.MulticastLocatorList.Locators {
		participant.networkFactory.FillDefaultLocatorPort(participant.DomaninID, &loc, participant.Att, true)
	}

	// Normalize unicast locators
	if att.UnicastLocatorList.Length() > 0 {
		participant.networkFactory.NormalizedLocators(&att.UnicastLocatorList)
	}
}

type createReaderCallback = func(guid *common.GUIDT, param *attributes.ReaderAttributes,
	persistenceSrv persistence.IPersistenceService, reliable bool) reader.IRTPSReader
type createWriterCallback = func(guid *common.GUIDT, param *attributes.WriterAttributes,
	persistenceSrv persistence.IPersistenceService, reliable bool) writer.IRTPSWriter

func (participant *RTPSParticipant) CreateWriter(param *attributes.WriterAttributes, payloadPool history.IPayloadPool,
	hist *history.WriterHistory, listen writer.IWriterListener,
	entityID *common.EntityIDT, isBuiltin bool) (bool, writer.IRTPSWriter) {
	if payloadPool == nil {
		log.Fatalln("Trying to create writer with null payload pool")
		return false, nil
	}

	callback := func(guid *common.GUIDT, pparam *attributes.WriterAttributes,
		persistenceSrv persistence.IPersistenceService, isReliable bool) writer.IRTPSWriter {
		if isReliable {
			if persistenceSrv != nil {
				log.Fatalln("notImpl")
				return nil
			} else {
				return writer.NewStatefulWriter(participant, guid, pparam,
					payloadPool, hist, listen)
			}
		} else {
			if persistenceSrv != nil {
				log.Fatalln("notImpl")
				return nil
			} else {
				poolConfig := history.FromHistoryAttributes(&hist.Att)
				cacheChange := history.NewCacheChangePool(poolConfig)
				return writer.NewStatelessWriter(participant, guid, param, payloadPool, cacheChange, hist, listen)
			}
		}
	}

	return participant.createWriter(param, entityID, isBuiltin, callback)
}

func (participant *RTPSParticipant) GetMinNetworkSendBufferSize() uint32 {
	return participant.networkFactory.GetMinSendBufferSize()
}

// Create non-existent SendResources based on the Locator list of the entity
func (participant *RTPSParticipant) createSendResources(pend endpoint.IEndpoint) bool {
	remoteLocators := &pend.GetAttributes().RemoteLocatorList
	if remoteLocators.Empty() {
		participant.networkFactory.GetDefaultOutputLocators(remoteLocators)
	}
	participant.sendResourceMutex.Lock()
	defer participant.sendResourceMutex.Unlock()
	for _, loc := range remoteLocators.Locators {
		if !participant.networkFactory.BuildSendResources(participant.sendResourceList, &loc) {
			log.Fatalln("Cannot create send resource for endpoint remote locator (", pend.GetGUID(),
				", ", loc, ")")
		}
	}

	return true
}

func (participant *RTPSParticipant) createWriter(param *attributes.WriterAttributes,
	entityID *common.EntityIDT, isBuiltin bool, callback createWriterCallback) (bool, writer.IRTPSWriter) {
	var reliabilityType string
	if param.EndpointAtt.ReliabilityKind == common.KReliable {
		reliabilityType = "RELIABLE"
	} else {
		reliabilityType = "BEST_EFFORT"
	}
	log.Println("Creating writer of type ", reliabilityType)
	ok, entID := participant.preprocessEndpointAttributes("writer", entityID, &param.EndpointAtt,
		common.KWriter, 0x03, 0x02)
	if !ok {
		return false, nil
	}
	throughputCtrl := &param.ThroughputController
	attThroughputCtrl := participant.Att.ThroughputController
	invalid := throughputCtrl.BytesPerPeriod != math.MaxUint32 && throughputCtrl.PeriodMillisecs != 0
	invalid = invalid || (attThroughputCtrl.BytesPerPeriod != math.MaxUint32 && attThroughputCtrl.PeriodMillisecs != 0)
	invalid = invalid && (param.PubMode != attributes.KAsynchronousWriter)
	if invalid {
		log.Fatalln("Writer has to be configured to publish asynchronously, because a flowcontroller was configured")
		return false, nil
	}

	// Special case for DiscoveryProtocol::BACKUP, which abuses persistence guid
	formerPersistenceGUID := param.EndpointAtt.PersistenceGUID
	if param.EndpointAtt.PersistenceGUID == common.KGuidUnknown {
		if participant.PersistenceGUID != common.KGuidUnknown {
			// Generate persistence guid from participant persistence guid
			param.EndpointAtt.PersistenceGUID.Prefix = participant.PersistenceGUID.Prefix
			param.EndpointAtt.PersistenceGUID.EntityID = *entID
		}
	}

	// Get persistence service
	ok, persistence := participant.getPersistenceService("writer", isBuiltin, &param.EndpointAtt)
	if !ok {
		return false, nil
	}
	participant.normalizeEndpointLocators(&param.EndpointAtt)
	guid := common.GUIDT{
		Prefix:   participant.GUID.Prefix,
		EntityID: *entID,
	}
	swriter := callback(&guid, param, persistence, param.EndpointAtt.ReliabilityKind == common.KReliable)

	// restore attributes
	param.EndpointAtt.PersistenceGUID = formerPersistenceGUID
	if swriter == nil {
		log.Fatalln("create writer failed")
		return false, nil
	}

	participant.createSendResources(swriter)
	if param.EndpointAtt.ReliabilityKind == common.KReliable {
		if !participant.createAndAssociateReceiverswithEndpoint(swriter) {
			log.Fatalln("createAndAssociateReceiverswithEndpoint failed")
			return false, nil
		}
	}

	participant.mutex.Lock()
	defer participant.mutex.Unlock()
	if isBuiltin {
		participant.asyncThread.Wakeup(swriter)
	} else {
		participant.UserWriterList = append(participant.UserWriterList, swriter)
	}

	// If the terminal throughput controller has proper user defined values, instantiate it
	throughputCtrl = &param.ThroughputController
	if throughputCtrl.BytesPerPeriod != math.MaxUint32 && throughputCtrl.PeriodMillisecs != 0 {
		controller := flowcontrol.NewThroughputController(&param.ThroughputController, swriter)
		swriter.AddFlowController(controller)
	}

	return true, swriter
}

func (participant *RTPSParticipant) createReader(param *attributes.ReaderAttributes,
	entityID *common.EntityIDT, isBuiltin, enable bool, callback createReaderCallback) (bool, reader.IRTPSReader) {

	var reliableType string
	if param.EndpointAtt.ReliabilityKind == common.KReliable {
		reliableType = "RELIABLE"
	} else {
		reliableType = "BEST_EFFORT"
	}
	log.Printf("Creating reader of type %v", reliableType)
	ok, entID := participant.preprocessEndpointAttributes("reader", entityID, &param.EndpointAtt, common.KReader, 0x04, 0x07)
	if !ok {
		return false, nil
	}

	// Special case for DiscoveryProtocol::BACKUP, which abuses persistence guid
	formerPersistenceGUID := param.EndpointAtt.PersistenceGUID
	if param.EndpointAtt.PersistenceGUID == common.KGuidUnknown {
		if participant.PersistenceGUID != common.KGuidUnknown {
			// Generate persistence guid from participant persistence guid
			param.EndpointAtt.PersistenceGUID = common.GUIDT{
				Prefix:   participant.PersistenceGUID.Prefix,
				EntityID: *entityID,
			}
		}
	}

	// Get persistence service
	ok, persistence := participant.getPersistenceService("reader", isBuiltin, &param.EndpointAtt)
	if !ok {
		return false, nil
	}

	participant.normalizeEndpointLocators(&param.EndpointAtt)
	guid := common.GUIDT{
		Prefix:   participant.GUID.Prefix,
		EntityID: *entID,
	}
	sReader := callback(&guid, param, persistence, param.EndpointAtt.ReliabilityKind == common.KReliable)

	// restore attributes
	param.EndpointAtt.PersistenceGUID = formerPersistenceGUID

	if sReader == nil {
		return false, nil
	}

	if param.EndpointAtt.ReliabilityKind == common.KReliable {
		participant.createSendResources(sReader)
	}

	if isBuiltin {
		trusted := trustedWriter(&sReader.GetGUID().EntityID)
		sReader.SetTrustedWriter(&trusted)
	}

	if enable {
		if !participant.createAndAssociateReceiverswithEndpoint(sReader) {
			return false, nil
		}
	}

	participant.mutex.Lock()
	defer participant.mutex.Unlock()
	participant.AllReaderList = append(participant.AllReaderList, sReader)
	if !isBuiltin {
		participant.UserReaderList = append(participant.UserReaderList, sReader)
	}

	return true, sReader
}

func (participant *RTPSParticipant) CreateReader(param *attributes.ReaderAttributes,
	payloadPool history.IPayloadPool, hist *history.ReaderHistory, listen reader.IReaderListener,
	entityID *common.EntityIDT, isBuiltin bool, enable bool) (bool, reader.IRTPSReader) {

	if payloadPool == nil {
		log.Fatalln("Trying to create reader with null payload pool")
		return false, nil
	}

	callback := func(guid *common.GUIDT, param *attributes.ReaderAttributes,
		persistenceServ persistence.IPersistenceService, isReliable bool) reader.IRTPSReader {

		if isReliable {
			if persistenceServ != nil {

			} else {
				return reader.NewStatefulReader(participant, guid, param, payloadPool, hist, listen)
			}

		} else {
			if persistenceServ != nil {

			} else {
				return reader.NewStatelessReader(participant, guid, param, payloadPool, hist, listen)
			}
		}

		return nil
	}

	return participant.createReader(param, entityID, isBuiltin, enable, callback)
}

func (participant *RTPSParticipant) GetEventResource() *resources.ResourceEvent {
	return participant.eventThr
}

// NewParticipant create new rtps participant
func NewParticipant(domainID uint32, pparam *attributes.RTPSParticipantAttributes, guidP,
	perstGUID *common.GUIDPrefixT, listen *RTPSParticipantListener) *RTPSParticipant {
	participant := RTPSParticipant{
		DomaninID:            domainID,
		Att:                  pparam,
		GUID:                 common.GUIDT{Prefix: *guidP, EntityID: *common.KEidRTPSParticipant},
		PersistenceGUID:      common.GUIDT{Prefix: *perstGUID, EntityID: *common.KEidRTPSParticipant},
		BuiltinProtocols:     nil,
		ResSemaphore:         utils.NewSemaphore(0),
		networkFactory:       network.NewNetworkFactory(),
		eventThr:             resources.NewResourceEvent(),
		IDCounter:            0,
		Listener:             listen,
		checkFn:              nil,
		IntraProcessOnly:     shouldBeIntraProcessOnly(pparam),
		hasShmTransport:      false,
		receiverResourceList: make(map[*ReceiverControlBlock]bool),
	}

	// Builtin transports by default
	if pparam.UseBuiltinTransports {
		descriptor := transport.NewUDPv4TransportDescriptor()
		descriptor.SendBufferSize = pparam.SendSocketBufferSize
		descriptor.RcvBufferSize = pparam.ListenSocketBufferSize
		participant.networkFactory.RegisterTransport(descriptor)

		if os.Getenv("SHM_TRANSPORT_BUILTIN") == "1" {
			log.Println("Use SHM_TRANSPORT_BUILTIN")

			shmTransport := transport.NewSharedMemTransportDescriptor()
			// We assume (Linux) UDP doubles the user socket buffer size in kernel, so
			// the equivalent segment size in SHM would be socket buffer size x 2
			var segmentSizeUDPequivalent uint32
			if participant.Att.SendSocketBufferSize >= participant.Att.ListenSocketBufferSize {
				segmentSizeUDPequivalent = participant.Att.SendSocketBufferSize * 2
			} else {
				segmentSizeUDPequivalent = participant.Att.ListenSocketBufferSize * 2
			}
			shmTransport.SegmentSize = segmentSizeUDPequivalent
			participant.hasShmTransport = participant.hasShmTransport || participant.networkFactory.RegisterTransport(shmTransport)
		}
	}

	if pparam.Builtin.DiscoveryConfig.Protocol == attributes.KDisPBackup {
		participant.PersistenceGUID = participant.GUID
	}

	protocol := pparam.Builtin.DiscoveryConfig.Protocol
	if protocol == attributes.KDisPBackup || protocol == attributes.KDisPClient ||
		protocol == attributes.KDisPServer {
		for _, descriptor := range pparam.UserTransports {
			if tcpDescriptor, ok := descriptor.(*transport.TCPTransportDescriptor); ok {
				if len(tcpDescriptor.ListeningPorts) == 0 {
					log.Fatal("Participant ", pparam.Name, " with GUID ", participant.GUID,
						" tries to use discovery server over TCP without providing a proper listening port")
				}
			}
		}
	}

	// User defined transports
	for _, descriptor := range pparam.UserTransports {
		var shmTransport transport.SharedMemTransportDescriptor
		if participant.networkFactory.RegisterTransport(descriptor) == true {
			if reflect.TypeOf(descriptor) == reflect.TypeOf(shmTransport) {
				participant.hasShmTransport = true
			}
		} else {
			// SHM transport could be disabled
			if reflect.TypeOf(descriptor) == reflect.TypeOf(shmTransport) {
				log.Fatal("Unable to Register SHM Transport. SHM Transport is not supported in the current platform.")
			} else {
				log.Fatal("User transport failed to register.")
			}
		}
	}

	participant.eventThr.InitThread()
	log.Println("finish participant eventThr initThread.")

	if participant.NetworkFactoryHasRegisteredTransports() == false {
		return nil
	}

	// Throughput controller, if the descriptor has valid values
	if pparam.ThroughputController.BytesPerPeriod != ^uint32(0) &&
		pparam.ThroughputController.PeriodMillisecs != 0 {
		controller := flowcontrol.NewThroughputController(pparam.ThroughputController, participant)
		participant.Controllers = append(participant.Controllers, controller)
	}

	// Creation of metatraffic locator and receiver resources
	multicastPort := participant.Att.Port.GetMulticastPort(participant.DomaninID)
	unicastPort := participant.Att.Port.GetUnicastPort(participant.DomaninID, uint32(participant.Att.ParticipantID))

	/* INSERT DEFAULT MANDATORY MULTICAST LOCATORS HERE */
	if participant.Att.Builtin.MetatrafficMulticastLocatorList.Length() == 0 &&
		participant.Att.Builtin.MetatrafficUnicastLocatorList.Length() == 0 {
		participant.networkFactory.GetDefaultMetatrafficMulticastLocators(participant.Att.Builtin.MetatrafficMulticastLocatorList, multicastPort)
		participant.networkFactory.NormalizedLocators(participant.Att.Builtin.MetatrafficMulticastLocatorList)

		participant.networkFactory.GetDefaultMetatrafficUnicastLocators(participant.Att.Builtin.MetatrafficUnicastLocatorList, unicastPort)
		participant.networkFactory.NormalizedLocators(participant.Att.Builtin.MetatrafficUnicastLocatorList)
	} else {
		for _, locator := range participant.Att.Builtin.MetatrafficMulticastLocatorList.Locators {
			participant.networkFactory.FillMetatrafficMulticastLocator(&locator, multicastPort)
		}
		participant.networkFactory.NormalizedLocators(participant.Att.Builtin.MetatrafficMulticastLocatorList)

		for _, locator := range participant.Att.Builtin.MetatrafficUnicastLocatorList.Locators {
			participant.networkFactory.FillMetatrafficUnicastLocator(&locator, multicastPort)
		}
		participant.networkFactory.NormalizedLocators(participant.Att.Builtin.MetatrafficUnicastLocatorList)
	}

	// Initial peers
	if participant.Att.Builtin.InitialPeersList.Empty() {
		participant.Att.Builtin.InitialPeersList = participant.Att.Builtin.MetatrafficMulticastLocatorList
	} else {
		for _, peer := range participant.Att.Builtin.InitialPeersList.Locators {
			participant.networkFactory.ConfigureInitialPeerLocator(participant.DomaninID, &peer, participant.Att)
		}
		participant.Att.Builtin.InitialPeersList = nil
	}

	// Creation of user locator and receiver resources
	hasLocatorsDefined := true
	//If no default locators are defined we define some.
	/* The reasoning here is the following.
	   If the parameters of the RTPS Participant don't hold default listening locators for the creation
	   of Endpoints, we make some for Unicast only.
	   If there is at least one listen locator of any kind, we do not create any default ones.
	   If there are no sending locators defined, we create default ones for the transports we implement.
	*/
	if participant.Att.DefaultUnicastLocatorList.Empty() &&
		participant.Att.DefaultMulticastLocatorList.Empty() {
		//Default Unicast Locators in case they have not been provided
		/* INSERT DEFAULT UNICAST LOCATORS FOR THE PARTICIPANT */
		hasLocatorsDefined = false
		participant.networkFactory.GetDefaultUnicastLocators(participant.DomaninID,
			participant.Att.DefaultUnicastLocatorList,
			participant.Att)
	} else {
		// Locator with port 0, calculate port.
		for _, loc := range participant.Att.DefaultUnicastLocatorList.Locators {
			participant.networkFactory.FillDefaultLocatorPort(participant.DomaninID, &loc, participant.Att, false)
		}
	}

	// Normalize unicast locators.
	participant.networkFactory.NormalizedLocators(participant.Att.DefaultUnicastLocatorList)

	if !hasLocatorsDefined {
		log.Printf("Created with NO default Unicast Locator List, adding Locators: %v",
			participant.Att.DefaultUnicastLocatorList)
	}

	if participant.IsIntraprocessOnly() {
		participant.Att.Builtin.MetatrafficUnicastLocatorList.Clear()
		participant.Att.DefaultUnicastLocatorList.Clear()
		participant.Att.DefaultMulticastLocatorList.Clear()
	}

	participant.createReceiverResources(participant.Att.Builtin.MetatrafficMulticastLocatorList, true, false)
	participant.createReceiverResources(participant.Att.Builtin.MetatrafficUnicastLocatorList, true, false)
	participant.createReceiverResources(participant.Att.DefaultUnicastLocatorList, true, false)
	participant.createReceiverResources(participant.Att.DefaultMulticastLocatorList, true, false)

	participant.asyncThread = writer.NewAsyncWriterThread()

	allowGrowingBuffers := participant.Att.Allocation.SendBuffers.Dynamic
	numSendBuffers := participant.Att.Allocation.SendBuffers.PreAllocatedNum
	if numSendBuffers == 0 {
		// Three buffers (user, events and async writer threads)
		numSendBuffers = 3
		// Add one buffer per reception thread
		numSendBuffers += uint64(len(participant.receiverResourceList))
	}

	// Create buffer pool
	participant.sendBuffers = message.NewSendBuffersManager(numSendBuffers, allowGrowingBuffers)
	participant.sendBuffers.Init(&participant.GUID.Prefix, participant.GetMaxMessageSize())

	participant.BuiltinProtocols = builtin.NewBuiltinProtocols()

	log.Printf("RTPSParticipant %v, with guidPrefix: %v", participant.Att.Name, participant.GUID.Prefix)

	return &participant
}
