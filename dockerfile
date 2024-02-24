FROM golang:1.21

ENV CGO_ENABLED 0
ENV GOOS "linux"
ENV GOOSE_DRIVER "postgres"
ENV GOOSE_DBSTRING "host=localhost user=postgres database=news_feed_bot password=postgres sslmode=disable"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download; go install github.com/pressly/goose/v3/cmd/goose@latest
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY config.hcl ./


RUN CGO_ENABLED=$CGO_ENABLED GOOS=$GOOS go build -o /val cmd/main.go

EXPOSE 8080

# Run
CMD ["/val"]