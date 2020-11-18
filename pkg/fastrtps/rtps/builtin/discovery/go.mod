module discovery

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/jiayin/Desktop/DDS/pkg/core/policy
	github.com/yeren0143/DDS/fastrtps/network => /home/jiayin/Desktop/DDS/pkg/fastrtps/network
	github.com/yeren0143/DDS/fastrtps/rtps/attributes => /home/jiayin/Desktop/DDS/pkg/fastrtps/rtps/attributes
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol => /home/jiayin/Desktop/DDS/pkg/fastrtps/rtps/flowcontrol
	github.com/yeren0143/DDS/fastrtps/rtps/resources => /home/jiayin/Desktop/DDS/pkg/fastrtps/rtps/resources
	github.com/yeren0143/DDS/fastrtps/transport => /home/jiayin/Desktop/DDS/pkg/fastrtps/transport
	github.com/yeren0143/DDS/fastrtps/utils => /home/jiayin/Desktop/DDS/pkg/fastrtps/utils
	github.com/yeren0143/DDS/types => /home/jiayin/Desktop/DDS/pkg/types
)

require (
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/attributes v0.0.0-00010101000000-000000000000
)
