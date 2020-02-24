VERSION :=$(shell cat version)

build:
	scripts/build.sh
format:
	scripts/exec_gofmt.sh
compile:
	GOOS=linux GOARCH=amd64 go build -o dist/terraform-provider-flowdock-linux_${VERSION}
	GOOS=darwin GOARCH=amd64 go build -o dist/terraform-provider-flowdock-darwin_${VERSION}
