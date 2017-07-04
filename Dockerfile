FROM golang:1.8

WORKDIR /go/src/app
COPY . .
COPY database database
COPY resolvers resolvers

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"] # ["app"]