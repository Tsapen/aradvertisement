# aradvertisement
Augmented reality advertisement project

Project for ar advertisement creation and management.  
***

## Private and public self-signed keys generation  
`openssl genrsa -out server.key 2048 && openssl ecparam -genkey -name secp384r1 -out server.key && openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`

## Running with default settings  
`cp config.example.json config.json && export ARA_CONFIG="$PWD/config.json" && go run ./internal/cmd/main.go `