module transport

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/pkg/common
	github.com/yeren0143/DDS/fastrtps/network => /home/jiayin/Desktop/DDS/pkg/fastrtps/network

)

require (
	github.com/yeren0143/DDS/common v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/network v0.0.0-00010101000000-000000000000
)
