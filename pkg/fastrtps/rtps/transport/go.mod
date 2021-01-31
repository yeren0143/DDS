module transport

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/jiayin/Desktop/DDS/DDS_go/pkg/common
	github.com/yeren0143/DDS/fastrtps/utils => /home/jiayin/Desktop/DDS/DDS_go/pkg/fastrtps/utils
)

require (
	github.com/yeren0143/DDS/common v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/utils v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68
)
