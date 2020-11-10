module participant

go 1.15

replace common => ../common
replace attributes => ../attributes
replace writer => ../writer
replace core/policy => ../../core/policy

require common v0.0.0-00010101000000-000000000000
