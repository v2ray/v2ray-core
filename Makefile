.DEFAULT_GOAL := build

GOFMT=gofmt
GC=go build

.PHONY: windows
windows:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 $(GC)  -o ./bin/v2ray ./main
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 $(GC)  -o ./bin/v2ctl ./infra/control/main/

.PHONY: linux
linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GC) -o ./bin/v2ray ./main
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GC) -o ./bin/v2ctl ./infra/control/main/

.PHONY: macos
macos:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GC) -o ./bin/v2ray ./main
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GC) -o ./bin/v2ctl ./infra/control/main/

.PHONY: clean
clean:
	rm ./bin/v2ctl
	rm ./bin/v2ray
