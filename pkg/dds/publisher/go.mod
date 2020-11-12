module publisher

go 1.15


replace (
	common => ../../common 
	core/policy => ../core/policy
	types => ../../types 
)