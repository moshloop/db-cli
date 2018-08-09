build:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	mkdir -p build
	gox -os="darwin linux windows" -arch="amd64"
	mkdir -p build/osx
	mkdir -p build/linux
	mkdir -p build/windows
	mv db-cli_darwin_amd64 build/osx/db-cli
	mv db-cli_linux_amd64 build/linux/db-cli
	mv db-cli_windows_amd64.exe build/windows/db-cli.exe
	cp README.md build/
	zip -r db-cli.zip build/*