FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64\
    ARA_CONFIG=config.docker.json

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build internal/cmd/main.go
    
# Export necessary port
EXPOSE 8000
EXPOSE 8001

# # Command to run when starting the container
CMD ["./main"]
