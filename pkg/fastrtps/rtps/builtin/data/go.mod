module data

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/DDS_go/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/policy
	github.com/yeren0143/DDS/fastrtps/message => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/message
	github.com/yeren0143/DDS/fastrtps/rtps/attributes => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/attributes
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/flowcontrol
	github.com/yeren0143/DDS/fastrtps/rtps/network => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/network
	github.com/yeren0143/DDS/fastrtps/rtps/qos => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/qos
	github.com/yeren0143/DDS/fastrtps/rtps/resources => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/resources
	github.com/yeren0143/DDS/fastrtps/rtps/transport => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/transport
	github.com/yeren0143/DDS/fastrtps/rtps/transport/shm => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/transport/shm
	github.com/yeren0143/DDS/fastrtps/utils => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/utils
	github.com/yeren0143/DDS/types => /home/jiayin/Desktop/DDS/DDS_go/pkg/types
)

require (
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/core/policy v0.0.0
	github.com/yeren0143/DDS/fastrtps/message v0.0.0-00010101000000-000000000000 // indirect
	github.com/yeren0143/DDS/fastrtps/rtps/attributes v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/network v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/qos v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/resources v0.0.0
)
