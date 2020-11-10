module rtps

go 1.15

replace (
	attributes => ./attributes
	common => ../../common
	flowcontrol => ./flowcontrol
	participant => ./participant
	qos => ./qos
	resources => ./resources
	transport => ../transport
	types => ../../types
	utils => ../utils
	writer => ./writer
	core/policy => ../../core/policy
)

require (
	attributes v0.0.0
	common v0.0.0
	participant v0.0.0-00010101000000-000000000000
	qos v0.0.0-00010101000000-000000000000 // indirect
	utils v0.0.0
	writer v0.0.0-00010101000000-000000000000 // indirect
)
