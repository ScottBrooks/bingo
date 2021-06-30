FROM golang:1.16 as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN GOOS=linux go build -o /go/bin/bingo ./cmd/bingo

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/bingo /bingo/bingo
COPY --from=build /go/src/app/assets /bingo/assets

WORKDIR /bingo
CMD ["/bingo/bingo"]
