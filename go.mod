module test

go 1.15

replace github.com/y/original-project => github.com/x/my-version v0.5.2

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-resty/resty/v2 v2.4.0
	github.com/inexio/go-monitoringplugin v0.0.0-20201117085742-ec06ef4904fe
	github.com/jessevdk/go-flags v1.4.0
	github.com/pkg/errors v0.9.1
)
