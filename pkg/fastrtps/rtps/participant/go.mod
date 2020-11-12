module participant

go 1.15

replace (
	common => ../../../common
	core/policy => ../../../core/policy
	dds/publisher/qos => ../../../dds/publisher/qos
	fastrtps/rtps/attributes => ../attributes
	fastrtps/rtps/builtin/discovery => ../builtin/discovery
	fastrtps/rtps/flowcontrol => ../flowcontrol
	fastrtps/rtps/qos => ../qos
	fastrtps/rtps/resources => ../resources
	fastrtps/rtps/writer => ../writer
    fastrtps/rtps/reader => ../reader
	fastrtps/transport => ../../transport
	fastrtps/utils => ../../utils
	types => ../../../types
    fastrtps/participant => ../../participant
)

require (
	common v0.0.0
	fastrtps/rtps/attributes v0.0.0-00010101000000-000000000000
	fastrtps/rtps/builtin/discovery v0.0.0-00010101000000-000000000000
	fastrtps/rtps/resources v0.0.0-00010101000000-000000000000
	fastrtps/rtps/writer v0.0.0-00010101000000-000000000000
)
