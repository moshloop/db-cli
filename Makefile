build:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	gox -os="darwin linux windows" -arch="amd64"
	mv db-cli_darwin_amd64 db-cli_osx
	mv db-cli_linux_amd64 db-cli
	mv db-cli_windows_amd64.exe db-cli.exe
	zip -r db-cli.zip README.md db-cli db-cli_osx db-cli.exe