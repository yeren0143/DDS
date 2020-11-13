module attributes

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/allride/Desktop/Ros/DDS/pkg/common
	github.com/yeren0143/DDS/core/policy => /home/allride/Desktop/Ros/DDS/pkg/core/policy
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol => /home/allride/Desktop/Ros/DDS/pkg/fastrtps/rtps/flowcontrol
	github.com/yeren0143/DDS/fastrtps/rtps/resources => /home/allride/Desktop/Ros/DDS/pkg/fastrtps/rtps/resources
	github.com/yeren0143/DDS/fastrtps/transport => /home/allride/Desktop/Ros/DDS/pkg/fastrtps/transport
	github.com/yeren0143/DDS/fastrtps/utils => /home/allride/Desktop/Ros/DDS/pkg/fastrtps/utils
	github.com/yeren0143/DDS/types => /home/allride/Desktop/Ros/DDS/pkg/types
)

require (
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/core/policy v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/flowcontrol v0.0.0
	github.com/yeren0143/DDS/fastrtps/rtps/resources v0.0.0
	github.com/yeren0143/DDS/fastrtps/transport v0.0.0
	github.com/yeren0143/DDS/fastrtps/utils v0.0.0
	github.com/yeren0143/DDS/types v0.0.0-00010101000000-000000000000 // indirect
)
