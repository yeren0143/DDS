module participant

go 1.15

replace common => ../common
replace fastrtps/rtps/attributes => ../attributes
replace fastrtps/rtps/writer => ../writer
replace core/policy => ../../core/policy
replace fastrtps/rtps/qos => ../qos
replace dds/publisher/qos => ../../../dds/publisher/qos

require common v0.0.0-00010101000000-000000000000
