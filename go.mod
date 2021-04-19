module github.com/moxspec/moxspec-metrics-agent

go 1.14

replace github.com/moxspec/moxspec-metrics-agent/promcli => ./promcli

require (
	github.com/digitalocean/go-smbios v0.0.0-20180907143718-390a4f403a8e // indirect
	github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/golang/snappy v0.0.3
	github.com/moxspec/moxspec v0.0.0-20210317194257-4bb7d102e8f3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/prometheus v1.8.2-0.20210419070145-a9a5f04ff9d6
	golang.org/x/crypto v0.0.0-20210415154028-4f45737414dc // indirect
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c // indirect
	golang.org/x/term v0.0.0-20210406210042-72f3dc4e9b72 // indirect
)
