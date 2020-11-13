module qos

go 1.15

replace github.com/yeren0143/DDS/core/policy => /home/allride/Desktop/Ros/DDS/pkg/core/policy

replace github.com/yeren0143/DDS/common => /home/allride/Desktop/Ros/DDS/pkg/common

replace github.com/yeren0143/DDS/types => /home/allride/Desktop/Ros/DDS/pkg/types

require (
	github.com/yeren0143/DDS/core/policy v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/types v0.0.0-00010101000000-000000000000 // indirect
)
