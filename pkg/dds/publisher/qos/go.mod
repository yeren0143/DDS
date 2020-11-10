module qos

go 1.15

replace (
	common => ../../../common //indirect
	core/policy => ../../../core/policy
	//types => ../../../types //indirect
    types => ${DDS_PATH}/pkg/types
)

require core/policy v0.0.0-00010101000000-000000000000
