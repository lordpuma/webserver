FROM golang:1.8

WORKDIR /go/src/app
COPY main.go main.go

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

EXPOSE 80

CMD ["go-wrapper", "run"] # ["app"]