FROM golang:1.21.6-alpine AS builder

WORKDIR /usr/src/app

RUN apk --no-cache add ca-certificates git && update-ca-certificates

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum .
RUN GOPROXY=direct GOFLAGS="-insecure" go mod download && go mod verify

COPY . .
RUN GOOS=linux GOARCH=arm GOPROXY=direct GOFLAGS="-insecure" go build -v -o /usr/local/bin/app main.go

#---------------------------------

From alpine

RUN apk --no-cache add curl

COPY --from=builder /usr/local/bin/app /usr/local/bin/app

CMD ["app"]
