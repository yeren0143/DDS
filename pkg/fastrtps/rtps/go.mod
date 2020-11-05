module rtps

go 1.15

replace common => ../common

replace utils => ../utils

require (
	common v0.0.0-00010101000000-000000000000
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	utils v0.0.0-00010101000000-000000000000
)
