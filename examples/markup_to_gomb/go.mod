module github.com/ernlel/gomb/examples/markup_to_gomb

go 1.25.0

require (
	github.com/ernlel/gomb v0.0.0
	github.com/ernlel/gomb/pkg/markup_to_gomb v0.0.0
)

require golang.org/x/net v0.56.0 // indirect

replace (
	github.com/ernlel/gomb => ../../
	github.com/ernlel/gomb/pkg/markup_to_gomb => ../../pkg/markup_to_gomb
)
