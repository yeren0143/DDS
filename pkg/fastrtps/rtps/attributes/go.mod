module attributes

go 1.15

replace common => ../common

replace flowcontrol => ../flowcontrol

replace transport => ../../transport

replace utils => ../../utils

replace resources => ../resources

require (
	common v0.0.0-00010101000000-000000000000
	flowcontrol v0.0.0-00010101000000-000000000000
	resources v0.0.0-00010101000000-000000000000
	transport v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)
