module participant

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/allride/Desktop/DDS/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/allride/Desktop/DDS/pkg/core/policy
	github.com/yeren0143/DDS/fastrtps/network => /home/allride/Desktop/DDS/pkg/fastrtps/network
	github.com/yeren0143/DDS/fastrtps/participant => /home/allride/Desktop/DDS/pkg/fastrtps/participant
	github.com/yeren0143/DDS/fastrtps/rtps/attributes => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/attributes
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/builtin/discovery
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/endpoint => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/builtin/discovery/endpoint
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/flowcontrol
	github.com/yeren0143/DDS/fastrtps/rtps/reader => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/reader
	github.com/yeren0143/DDS/fastrtps/rtps/resources => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/resources
	github.com/yeren0143/DDS/fastrtps/rtps/writer => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/writer
	github.com/yeren0143/DDS/fastrtps/transport => /home/allride/Desktop/DDS/pkg/fastrtps/transport
	github.com/yeren0143/DDS/fastrtps/utils => /home/allride/Desktop/DDS/pkg/fastrtps/utils
	github.com/yeren0143/DDS/types => /home/allride/Desktop/DDS/pkg/types
)

require (
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/attributes v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/reader v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/resources v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/writer v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/utils v0.0.0
)
