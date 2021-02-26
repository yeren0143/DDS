module participant

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/DDS_go/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/policy
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/edp => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/builtin/discovery/edp
	github.com/yeren0143/DDS/fastrtps/rtps/reader => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/reader
	github.com/yeren0143/DDS/fastrtps/rtps/writer => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/writer
	github.com/yeren0143/DDS/types => /home/jiayin/Desktop/DDS/DDS_go/pkg/types

)