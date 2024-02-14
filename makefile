help:
	@cat ./makefile

all:
	@termline grey
	GOOS=linux GOARCH=amd64 go build -o build/ironclad-linux-amd64/ironclad ./ironclad
	GOOS=darwin GOARCH=amd64 go build -o build/ironclad-mac-amd64/ironclad ./ironclad
	GOOS=darwin GOARCH=arm64 go build -o build/ironclad-mac-arm64/ironclad ./ironclad
	GOOS=windows GOARCH=amd64 go build -o build/ironclad-windows-amd64/ironclad.exe ./ironclad
	@termline grey
	@tree build
	@termline grey
	@shasum -a 256 build/*/*
	@termline grey
	@mkdir -p zipped
	@cd build && zip -r ../zipped/ironclad-linux-amd64.zip ironclad-linux-amd64 > /dev/null
	@cd build && zip -r ../zipped/ironclad-mac-amd64.zip ironclad-mac-amd64 > /dev/null
	@cd build && zip -r ../zipped/ironclad-mac-arm64.zip ironclad-mac-arm64 > /dev/null
	@cd build && zip -r ../zipped/ironclad-windows-amd64.zip ironclad-windows-amd64 > /dev/null
	@tree zipped
	@termline grey

clean:
	rm -rf ./build
	rm -rf ./zipped
