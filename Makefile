setup:
	mkdir -p {dist/linux,dist/windows,dist,macos}

build-windows: dependencies setup
	export GOOS=windows; go build main.go; cp main dist/windows/user-api

build-macos: dependencies setup
	export GOOS=darwin; go build main.go; cp main dist/macos/user-api

build-linux: dependencies setup
	export GOOS=linux; go build main.go; cp main dist/linux/user-api

dependencies:
	go mod vendor

build-container-image: build-linux
	cp dist/linux/user-api docker/;
	cd docker; docker build -t udaykiranr/user-api:`cat ../.version` .;
	docker push udaykiranr/user-api:`cat .version`;

clean:
	rm main; rm docker/user-api; rm -rf dist