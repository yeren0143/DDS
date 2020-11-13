module rtps

go 1.15

replace (
	common => ../../common
	core/policy => ../../core/policy
	dds/publisher/qos => ../../dds/publisher/qos
	fastrtps/participant => ../participant
	fastrtps/rtps/attributes => ./attributes
	fastrtps/rtps/builtin/discovery => ./builtin/discovery
	fastrtps/rtps/builtin/discovery/endpoint => ./builtin/discovery/endpoint
	fastrtps/rtps/flowcontrol => ./flowcontrol
	fastrtps/rtps/participant => ./participant
	fastrtps/rtps/qos => ./qos
	fastrtps/rtps/reader => ./reader
	fastrtps/rtps/resources => ./resources
	fastrtps/rtps/writer => ./writer
	fastrtps/transport => ../transport
	fastrtps/utils => ../utils
	types => ../../types
)

require (
	common v0.0.0
	core/policy v0.0.0 // indirect
	fastrtps/participant v0.0.0 // indirect
	fastrtps/rtps/attributes v0.0.0
	fastrtps/rtps/participant v0.0.0
	fastrtps/rtps/reader v0.0.0 // indirect
	fastrtps/rtps/resources v0.0.0
	fastrtps/utils v0.0.0
)
