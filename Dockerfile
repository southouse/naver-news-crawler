FROM golang:1.21.0-alpine3.18 AS gobuilder

WORKDIR /
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o naver-news-crawler main.go

CMD ["/naver-news-crawler"]