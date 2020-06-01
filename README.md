# aradvertisement
Augmented reality advertisement project

Project for ar advertisement creation and management provides funtional of adding objects availible for everybody. 

## Requirements:
`docker v19.03`

or  

`go v1.13`  
`redis v5.0`  
`postgres v10.12`  
`npm v3.5`  
***
## Installation

`git clone github.com/Tsapen/aradvertisement`  
`cd ./aradvertisement`  

`## Private and public self-signed keys generation`  
`openssl genrsa -out server.key 2048 && openssl ecparam -genkey -name secp384r1 -out server.key && openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`

`mv server.crt server.key certs`

## Running with default settings on local machine
`export ARA_CONFIG="$PWD/config.json" && go run ./internal/cmd/main.go `

## Running in docker
`docker-compose up`

## Testing
`export ARA_TEST_CONFIG="$PWD/config_test.json" &&  go test -v ./internal/cmd/`