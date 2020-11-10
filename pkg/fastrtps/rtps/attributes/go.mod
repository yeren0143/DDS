module attributes

go 1.15

replace common => ../../../common

replace fastrtps/transport => ../../transport

replace fastrtps/utils => ../../utils

replace fastrtps/rtps/resources => ../resources

replace core/policy => ../../../core/policy

replace types => ../../../types // indirect

replace fastrtps/rtps/flowcontrol => ../flowcontrol

require (
	common v0.0.0
	core/policy v0.0.0-00010101000000-000000000000
	fastrtps/rtps/flowcontrol v0.0.0-00010101000000-000000000000
	fastrtps/rtps/resources v0.0.0-00010101000000-000000000000
	fastrtps/transport v0.0.0-00010101000000-000000000000
	fastrtps/utils v0.0.0-00010101000000-000000000000
)
