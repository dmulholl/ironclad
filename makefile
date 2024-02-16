help:
	@cat ./makefile

all: clean
	@printf "\e[1;32mCompiling\e[0m ...\n"
	GOOS=linux GOARCH=amd64 go build -o build/ironclad-linux-amd64/ironclad ./ironclad
	GOOS=darwin GOARCH=amd64 go build -o build/ironclad-mac-amd64/ironclad ./ironclad
	GOOS=darwin GOARCH=arm64 go build -o build/ironclad-mac-arm64/ironclad ./ironclad
	GOOS=windows GOARCH=amd64 go build -o build/ironclad-windows-amd64/ironclad.exe ./ironclad
	@printf "\n\e[1;32mBinaries\e[0m ...\n"
	@tree build
	@printf "\n\e[1;32mSHA 256 checksums\e[0m ...\n"
	@shasum -a 256 build/*/*
	@printf "\n\e[1;32mZipping\e[0m ...\n"
	@mkdir -p build/zipped
	@cd build && zip -r ./zipped/ironclad-linux-amd64.zip ironclad-linux-amd64 > /dev/null
	@cd build && zip -r ./zipped/ironclad-mac-amd64.zip ironclad-mac-amd64 > /dev/null
	@cd build && zip -r ./zipped/ironclad-mac-arm64.zip ironclad-mac-arm64 > /dev/null
	@cd build && zip -r ./zipped/ironclad-windows-amd64.zip ironclad-windows-amd64 > /dev/null
	@tree build/zipped

clean:
	@rm -rf ./build
