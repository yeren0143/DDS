module policy

go 1.15

replace (
	github.com/yeren0143/DDS/common => /home/allride/Desktop/DDS/pkg/common
	github.com/yeren0143/DDS/types => /home/allride/Desktop/DDS/pkg/types
)

require (
	github.com/yeren0143/DDS/common v0.0.0
	github.com/yeren0143/DDS/types v0.0.0-00010101000000-000000000000
)
