module qos

go 1.15

replace (
	common => ../../../common //indirect
	core/policy => ../../core/policy
	types => ../../../types //indirect
)

require core/policy v0.0.0