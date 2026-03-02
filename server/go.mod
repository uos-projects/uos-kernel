module github.com/uos-projects/uos-kernel/server

go 1.23

toolchain go1.24.11

require (
	github.com/uos-projects/uos-kernel/actor v0.0.0
	github.com/uos-projects/uos-kernel/kernel v0.0.0
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/uos-projects/uos-kernel/meta v0.0.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/uos-projects/uos-kernel/actor => ../actor

replace github.com/uos-projects/uos-kernel/kernel => ../kernel

replace github.com/uos-projects/uos-kernel/meta => ../meta
