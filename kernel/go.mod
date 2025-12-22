module github.com/uos-projects/uos-kernel/kernel

go 1.23

toolchain go1.24.11

require (
	github.com/uos-projects/uos-kernel/actors v0.0.0
	github.com/uos-projects/uos-kernel/meta v0.0.0
)

require (
	github.com/apache/thrift v0.22.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/uos-projects/uos-kernel/actors => ../actors

replace github.com/uos-projects/uos-kernel/meta => ../meta
