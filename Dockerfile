FROM golang:1-alpine
WORKDIR /
COPY shame.go /shame.go
RUN go build -o /entrypoint
CMD /entrypoint