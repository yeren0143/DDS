module protocol

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/DDS_go/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/policy
	github.com/yeren0143/DDS/core/status => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/status
	github.com/yeren0143/DDS/fastrtps/message => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/message
	github.com/yeren0143/DDS/fastrtps/rtps/attributes => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/attributes
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/data => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/builtin/data
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/edp => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/builtin/discovery/edp
	github.com/yeren0143/DDS/fastrtps/rtps/endpoint => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/endpoint
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/flowcontrol
	github.com/yeren0143/DDS/fastrtps/rtps/history => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/history
	github.com/yeren0143/DDS/fastrtps/rtps/network => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/network
	github.com/yeren0143/DDS/fastrtps/rtps/qos => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/qos
	github.com/yeren0143/DDS/fastrtps/rtps/reader => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/reader
	github.com/yeren0143/DDS/fastrtps/rtps/resources => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/resources
	github.com/yeren0143/DDS/fastrtps/rtps/transport => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/transport
	github.com/yeren0143/DDS/fastrtps/rtps/writer => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/writer
	github.com/yeren0143/DDS/fastrtps/utils => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/utils
	github.com/yeren0143/DDS/types => /home/jiayin/Desktop/DDS/DDS_go/pkg/types
)

require (
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/attributes v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/data v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/history v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/network v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/reader v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/resources v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/writer v0.0.0-00010101000000-000000000000
)
