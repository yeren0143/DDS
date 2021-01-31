module message

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/DDS_go/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/policy
	github.com/yeren0143/DDS/fastrtps/rtps/endpoint => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/endpoint
	github.com/yeren0143/DDS/fastrtps/rtps/reader => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/reader
	github.com/yeren0143/DDS/fastrtps/rtps/writer => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/writer
	github.com/yeren0143/DDS/types => /home/jiayin/Desktop/DDS/DDS_go/pkg/types
)

require (
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/core/policy v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/reader v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/writer v0.0.0-00010101000000-000000000000
)
