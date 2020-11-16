module transport

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/allride/Desktop/DDS/pkg/common
	github.com/yeren0143/DDS/fastrtps/network => /home/allride/Desktop/DDS/pkg/fastrtps/network
	github.com/yeren0143/DDS/fastrtps/utils => /home/allride/Desktop/DDS/pkg/fastrtps/utils

)

require (
	github.com/yeren0143/DDS/common v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/network v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/utils v0.0.0-00010101000000-000000000000
)
