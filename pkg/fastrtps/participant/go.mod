module participant

go 1.15

replace (
	github.com/yeren0143/DDS/fastrtps/rtps/builtin/discovery/endpoint => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/builtin/discovery/endpoint
	github.com/yeren0143/DDS/fastrtps/rtps/reader => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/reader
	github.com/yeren0143/DDS/fastrtps/rtps/writer => /home/allride/Desktop/DDS/pkg/fastrtps/rtps/writer

)

require (
	github.com/yeren0143/DDS/fastrtps/rtps/reader v0.0.0-00010101000000-000000000000
	github.com/yeren0143/DDS/fastrtps/rtps/writer v0.0.0-00010101000000-000000000000
)
