build:
	$(MAKE) build-windows
	$(MAKE) build-linux
	$(MAKE) build-arm

build-windows:
	$(eval export GOOS=windows)
	$(eval export GOARCH=amd64)
	go build -o bin/gmf-windows-amd64.exe .

build-linux:
	$(eval export GOOS=linux)
	$(eval export GOARCH=amd64)
	go build -o bin/gmf-linux-amd64 .

build-arm:
	$(eval export GOOS=linux)
	$(eval export GOARCH=arm)
	$(eval export GOARM=7)
	go build -o bin/gmf-linux-armv7 .
