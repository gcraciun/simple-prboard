FROM golang:1.19 AS builder

WORKDIR /app

COPY . .

RUN go mod download

# To build from UBUNTU
#RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/simple-prboard
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/simple-prboard

# To build from UBUNTU
#FROM ubuntu:22.04
FROM scratch

WORKDIR /app

# To build from UBUNTU
# Install the ca-certificates package
#RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# To build from SCRATCH
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/httpRoot/ /app/httpRoot/
COPY --from=builder /app/templates/ /app/templates/
COPY --from=builder /go/bin/simple-prboard /app/

ENTRYPOINT ["/app/simple-prboard"]

