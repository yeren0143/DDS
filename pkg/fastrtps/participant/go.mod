module participant

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/DDS_go/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/policy
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/endpoint => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/builtin/discovery/endpoint
	github.com/yeren0143/DDS/fastrtps/rtps/reader => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/reader
	github.com/yeren0143/DDS/fastrtps/rtps/writer => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/writer
	github.com/yeren0143/DDS/types => /home/jiayin/Desktop/DDS/DDS_go/pkg/types

)

require (
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/endpoint v0.0.0-00010101000000-000000000000 // indirect
	github.com/yeren0143/DDS/fastrtps/rtps/reader v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/writer v0.0.0-00010101000000-000000000000
)
