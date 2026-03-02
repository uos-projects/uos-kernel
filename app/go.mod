module github.com/uos-projects/uos-kernel/app

go 1.23

toolchain go1.24.11

require github.com/uos-projects/uos-kernel/kernel v0.0.0

require (
	github.com/uos-projects/uos-kernel/actor v0.0.0 // indirect
	github.com/uos-projects/uos-kernel/meta v0.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/uos-projects/uos-kernel/kernel => ../kernel

replace github.com/uos-projects/uos-kernel/actor => ../actor

replace github.com/uos-projects/uos-kernel/meta => ../meta
