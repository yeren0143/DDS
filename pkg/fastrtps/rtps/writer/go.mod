module writer

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/DDS_go/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/policy
	github.com/yeren0143/DDS/core/status => /home/jiayin/Desktop/DDS/DDS_go/pkg/core/status
	github.com/yeren0143/DDS/fastrtps/rtps/attributes => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/rtps/attributes
)

require github.com/yeren0143/DDS/common v0.0.0-00010101000000-000000000000
