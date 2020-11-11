module rtps

go 1.15

replace (
	common => ../../common
	core/policy => ../../core/policy
	dds/publisher/qos => ../../dds/publisher/qos
	fastrtps/rtps/attributes => ./attributes
	fastrtps/rtps/flowcontrol => ./flowcontrol
	fastrtps/rtps/participant => ./participant
	fastrtps/rtps/qos => ./qos
	fastrtps/rtps/resources => ./resources
	fastrtps/rtps/writer => ./writer
	fastrtps/transport => ../transport
	fastrtps/utils => ../utils
	types => ../../types
)

require (
	common v0.0.0
	dds/publisher/qos v0.0.0-00010101000000-000000000000
	fastrtps/rtps/attributes v0.0.0
	fastrtps/rtps/participant v0.0.0-00010101000000-000000000000
	fastrtps/rtps/qos v0.0.0-00010101000000-000000000000 // indirect
	fastrtps/rtps/writer v0.0.0-00010101000000-000000000000 // indirect
	fastrtps/utils v0.0.0
)
