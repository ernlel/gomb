module github.com/ernlel/gomb/examples/html-elements

go 1.23.1

require (
	github.com/ernlel/gomb v0.0.0
	github.com/ernlel/gomb/pkg/html v0.0.0
)

replace (
	github.com/ernlel/gomb => ../../
	github.com/ernlel/gomb/pkg/html => ../../pkg/html
)
