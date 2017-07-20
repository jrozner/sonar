FROM golang:1.8 as build-stage
WORKDIR /go/src/sonar
COPY cmd/sonar/main.go .
RUN go get && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sonar .

FROM alpine:latest
WORKDIR /app
COPY --from=build-stage /go/src/sonar /app
COPY wordlist.txt /app
ENTRYPOINT ["/app/sonar"]
