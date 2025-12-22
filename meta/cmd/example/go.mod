module example

go 1.22.2

require (
	github.com/uos-projects/uos-kernel/actors v0.0.0
	github.com/uos-projects/uos-kernel/kernel v0.0.0
)

require (
	github.com/uos-projects/uos-kernel/meta v0.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/uos-projects/uos-kernel/kernel => ../../../kernel

replace github.com/uos-projects/uos-kernel/actors => ../../../actors

replace github.com/uos-projects/uos-kernel/meta => ../../../meta
