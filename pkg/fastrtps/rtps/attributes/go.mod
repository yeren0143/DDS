module attributes

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/allride/Desktop/DDS/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/allride/Desktop/DDS/pkg/core/policy
	github.com/yeren0143/DDS/fastrtps/network => /home/allride/Desktop/DDS/pkg/fastrtps/network
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/flowcontrol
	github.com/yeren0143/DDS/fastrtps/rtps/resources => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/resources
	github.com/yeren0143/DDS/fastrtps/transport => /home/allride/Desktop/DDS/pkg/fastrtps/transport
	github.com/yeren0143/DDS/fastrtps/utils => /home/allride/Desktop/DDS/pkg/fastrtps/utils
	github.com/yeren0143/DDS/types => /home/allride/Desktop/DDS/pkg/types
)

require (
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/core/policy v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/resources v0.0.0
	github.com/yeren0143/DDS/fastrtps/utils v0.0.0
)
