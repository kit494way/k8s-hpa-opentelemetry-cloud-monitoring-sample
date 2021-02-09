FROM golang:1.15-alpine
WORKDIR /go/src/app
COPY . .
RUN go install -v ./main.go
CMD ["main"]
