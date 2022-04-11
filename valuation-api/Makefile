.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/makes makes/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/models models/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/vin vin/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/trade-price trade-price/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
