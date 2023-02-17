all: build

APP_NAME = vscode.ext
PKGDIR_TMP = ${TMPDIR}golang

.pre-build:
	mkdir -p build

init:
	go mod init github.com/nachorpaez/osquery-vscode-extension

download:
	go mod download

clean:
	rm -rf build/
	rm -rf ${PKGDIR_TMP}_darwin

test:
	go test -v ./... 

build: .pre-build
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/${APP_NAME}-amd64 -pkgdir ${PKGDIR_TMP}
	GOOS=darwin GOARCH=arm64 go build -o build/darwin/${APP_NAME}-arm64 -pkgdir ${PKGDIR_TMP}
	lipo -create -output build/darwin/${APP_NAME} build/darwin/${APP_NAME}-amd64 build/darwin/${APP_NAME}-arm64
	GOOS=windows GOARCH=amd64  go build -o build/windows/${APP_NAME}.amd64.ext
	GOOS=windows GOARCH=arm64  go build -o build/windows/${APP_NAME}.arm64.ext
	GOOS=linux GOARCH=amd64  go build -o build/linux/${APP_NAME}.amd64.ext
	/bin/rm build/darwin/${APP_NAME}-amd64
	/bin/rm build/darwin/${APP_NAME}-arm64

osqueryi: build
	sleep 2
	osqueryi --extension=build/darwin/vscode.ext --allow_unsafe
