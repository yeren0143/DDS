package participant

import (
	"github.com/yeren0143/DDS/common"
	"github.com/yeren0143/DDS/fastrtps/message"
	"github.com/yeren0143/DDS/fastrtps/rtps/attributes"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin"
	"github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/participant"
	"github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol"
	"github.com/yeren0143/DDS/fastrtps/rtps/network"
	"github.com/yeren0143/DDS/fastrtps/rtps/reader"
	"github.com/yeren0143/DDS/fastrtps/rtps/resources"
	"github.com/yeren0143/DDS/fastrtps/rtps/transport"
	"github.com/yeren0143/DDS/fastrtps/rtps/writer"
	"github.com/yeren0143/DDS/fastrtps/utils"
	"log"
	"os"
	"reflect"
	"sync"
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

var _ builtin.IProtocolUser = (*RTPSParticipant)(nil)
var _ participant.IRTPSParticipant = (*RTPSParticipant)(nil)

//RTPSParticipant allows the creation and removal of writers and readers. It manages the send and receive threads.
type RTPSParticipant struct {
	DomaninID                 uint32
	Att                       *attributes.RTPSParticipantAttributes
	GUID                      common.GUIDT
	PersistenceGUID           common.GUIDT // Persistence guid of the RTPSParticipant
	EventThr                  *resources.ResourceEvent
	BuiltinProtocols          *builtin.Protocols
	ResSemaphore              *utils.Semaphore // Semaphore to wait for the listen thread creation.
	IDCounter                 uint32           // Id counter to correctly assign the ids to writers and readers.
	AllWriterList             []writer.IRTPSWriter
	AllReaderList             []reader.IRTPSReader
	Controllers               []flowcontrol.IFlowController
	networkFactory            *network.NetFactory
	Listener                  *RTPSParticipantListener
	IntraProcessOnly          bool
	hasShmTransport           bool
	checkFn                   typeCheckFn
	receiverResourceList      map[*ReceiverControlBlock]bool
	sendBuffers               *message.SendBuffersManager
	receiverResourceListMutex sync.Mutex
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

func (participant *RTPSParticipant) NetFactory() *network.NetFactory {
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
	if participant.BuiltinProtocols.InitBuiltinProtocol(participant, participant.Att.Builtin) == false {
		log.Fatal("The builtin protocols were not correctly initialized")
	}

	for receiver := range participant.receiverResourceList {
		receiver.ReceiverRes.RegisterReceiver(receiver.MsgReceiver)
	}
}

func (participant *RTPSParticipant) IsIntraprocessOnly() bool {
	return participant.IntraProcessOnly
}

func (participant *RTPSParticipant) applyLocatorAdaptRule(loc *common.Locator) *common.Locator {
	// This is a completely made up rule
	// It is transport responsability to interpret this new port.
	loc.Port += uint32(participant.Att.Port.ParticipantIDGain)
	return loc
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
		EventThr:             resources.NewResourceEvent(),
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

	participant.EventThr.InitThread()

	if participant.NetworkFactoryHasRegisteredTransports() == false {
		return nil
	}

	// Throughput controller, if the descriptor has valid values
	if pparam.ThroghputController.BytesPerPeriod != ^uint32(0) &&
		pparam.ThroghputController.PeriodMillisecs != 0 {
		controller := flowcontrol.NewThroughputController(pparam.ThroghputController, participant)
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
