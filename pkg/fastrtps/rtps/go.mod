module rtps

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/DDS_go/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/policy
	github.com/yeren0143/DDS/core/status => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/status
	github.com/yeren0143/DDS/dds/publisher/qos => /home/jiayin/Desktop/DDS/DDS_go/pkg/dds/publisher/qos
	github.com/yeren0143/DDS/fastrtps/message => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/message
	github.com/yeren0143/DDS/fastrtps/participant => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/participant
	github.com/yeren0143/DDS/fastrtps/rtps/attributes => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/attributes
	github.com/yeren0143/DDS/fastrtps/rtps/builtin => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/builtin
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/data => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/builtin/data
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/endpoint => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/builtin/discovery/endpoint
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/participant => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/builtin/discovery/participant
	github.com/yeren0143/DDS/fastrtps/rtps/endpoint => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/endpoint
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/flowcontrol
	github.com/yeren0143/DDS/fastrtps/rtps/history => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/history
	github.com/yeren0143/DDS/fastrtps/rtps/network => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/network
	github.com/yeren0143/DDS/fastrtps/rtps/participant => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/participant
	github.com/yeren0143/DDS/fastrtps/rtps/persistence => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/persistence
	github.com/yeren0143/DDS/fastrtps/rtps/qos => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/qos
	github.com/yeren0143/DDS/fastrtps/rtps/reader => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/reader
	github.com/yeren0143/DDS/fastrtps/rtps/resources => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/resources
	github.com/yeren0143/DDS/fastrtps/rtps/transport => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/transport
	github.com/yeren0143/DDS/fastrtps/rtps/writer => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/writer
	github.com/yeren0143/DDS/fastrtps/utils => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/utils
	github.com/yeren0143/DDS/types => /home/jiayin/Desktop/DDS/DDS_go/pkg/types
)

require (
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/attributes v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/data v0.0.0-00010101000000-000000000000 // indirect
	github.com/yeren0143/DDS/fastrtps/rtps/participant v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/persistence v0.0.0-00010101000000-000000000000 // indirect
	github.com/yeren0143/DDS/fastrtps/utils v0.0.0
)
