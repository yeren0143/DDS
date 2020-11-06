module policy

go 1.15

replace common => ../../common

replace types => ../../types

require (
	common v0.0.0-00010101000000-000000000000
	types v0.0.0-00010101000000-000000000000
)
