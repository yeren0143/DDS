package policy

type ParameterId_t uint16

const (
	PID_PAD                                 ParameterId_t = 0x0000
	PID_SENTINEL                            ParameterId_t = 0x0001
	PID_USER_DATA                           ParameterId_t = 0x002c
	PID_TOPIC_NAME                          ParameterId_t = 0x0005
	PID_TYPE_NAME                           ParameterId_t = 0x0007
	PID_GROUP_DATA                          ParameterId_t = 0x002d
	PID_TOPIC_DATA                          ParameterId_t = 0x002e
	PID_DURABILITY                          ParameterId_t = 0x001d
	PID_DURABILITY_SERVICE                  ParameterId_t = 0x001e
	PID_DEADLINE                            ParameterId_t = 0x0023
	PID_LATENCY_BUDGET                      ParameterId_t = 0x0027
	PID_LIVELINESS                          ParameterId_t = 0x001b
	PID_RELIABILITY                         ParameterId_t = 0x001A
	PID_LIFESPAN                            ParameterId_t = 0x002b
	PID_DESTINATION_ORDER                   ParameterId_t = 0x0025
	PID_HISTORY                             ParameterId_t = 0x0040
	PID_RESOURCE_LIMITS                     ParameterId_t = 0x0041
	PID_OWNERSHIP                           ParameterId_t = 0x001f
	PID_OWNERSHIP_STRENGTH                  ParameterId_t = 0x0006
	PID_PRESENTATION                        ParameterId_t = 0x0021
	PID_PARTITION                           ParameterId_t = 0x0029
	PID_TIME_BASED_FILTER                   ParameterId_t = 0x0004
	PID_TRANSPORT_PRIORITY                  ParameterId_t = 0x0049
	PID_PROTOCOL_VERSION                    ParameterId_t = 0x0015
	PID_VendorIDT                           ParameterId_t = 0x0016
	PID_UNICAST_LOCATOR                     ParameterId_t = 0x002f
	PID_MULTICAST_LOCATOR                   ParameterId_t = 0x0030
	PID_MULTICAST_IPADDRESS                 ParameterId_t = 0x0011
	PID_DEFAULT_UNICAST_LOCATOR             ParameterId_t = 0x0031
	PID_DEFAULT_MULTICAST_LOCATOR           ParameterId_t = 0x0048
	PID_METATRAFFIC_UNICAST_LOCATOR         ParameterId_t = 0x0032
	PID_METATRAFFIC_MULTICAST_LOCATOR       ParameterId_t = 0x0033
	PID_DEFAULT_UNICAST_IPADDRESS           ParameterId_t = 0x000c
	PID_DEFAULT_UNICAST_PORT                ParameterId_t = 0x000e
	PID_METATRAFFIC_UNICAST_IPADDRESS       ParameterId_t = 0x0045
	PID_METATRAFFIC_UNICAST_PORT            ParameterId_t = 0x000d
	PID_METATRAFFIC_MULTICAST_IPADDRESS     ParameterId_t = 0x000b
	PID_METATRAFFIC_MULTICAST_PORT          ParameterId_t = 0x0046
	PID_EXPECTS_INLINE_QOS                  ParameterId_t = 0x0043
	PID_PARTICIPANT_MANUAL_LIVELINESS_COUNT ParameterId_t = 0x0034
	PID_PARTICIPANT_BUILTIN_ENDPOINTS       ParameterId_t = 0x0044
	PID_PARTICIPANT_LEASE_DURATION          ParameterId_t = 0x0002
	PID_CONTENT_FILTER_PROPERTY             ParameterId_t = 0x0035
	PID_PARTICIPANT_GUID                    ParameterId_t = 0x0050
	PID_PARTICIPANT_ENTITYID                ParameterId_t = 0x0051
	PID_GROUP_GUID                          ParameterId_t = 0x0052
	PID_GROUP_ENTITYID                      ParameterId_t = 0x0053
	PID_BUILTIN_ENDPOINT_SET                ParameterId_t = 0x0058
	PID_PROPERTY_LIST                       ParameterId_t = 0x0059
	PID_TYPE_MAX_SIZE_SERIALIZED            ParameterId_t = 0x0060
	PID_ENTITY_NAME                         ParameterId_t = 0x0062
	PID_TYPE_IDV1                           ParameterId_t = 0x0069
	PID_KEY_HASH                            ParameterId_t = 0x0070
	PID_STATUS_INFO                         ParameterId_t = 0x0071
	PID_TYPE_OBJECTV1                       ParameterId_t = 0x0072
	PID_ENDPOINT_GUID                       ParameterId_t = 0x005a
	PID_IDENTITY_TOKEN                      ParameterId_t = 0x1001
	PID_PERMISSIONS_TOKEN                   ParameterId_t = 0x1002
	PID_DATA_TAGS                           ParameterId_t = 0x1003
	PID_ENDPOINT_SECURITY_INFO              ParameterId_t = 0x1004
	PID_PARTICIPANT_SECURITY_INFO           ParameterId_t = 0x1005
	PID_IDENTITY_STATUS_TOKEN               ParameterId_t = 0x1006
	PID_PERSISTENCE_GUID                    ParameterId_t = 0x8002
	PID_RELATED_SAMPLE_IDENTITY             ParameterId_t = 0x800f
	PID_DATA_REPRESENTATION                 ParameterId_t = 0x0073
	PID_TYPE_CONSISTENCY_ENFORCEMENT        ParameterId_t = 0x0074
	PID_TYPE_INFORMATION                    ParameterId_t = 0x0075
	PID_DISABLE_POSITIVE_ACKS               ParameterId_t = 0x8005
	//PID_RELATED_SAMPLE_IDENTITY = 0x0083
)

type Parameter_t struct {
	Pid    ParameterId_t
	Length uint16
}

type ParameterKey_t struct {
	Parameter_t
}
