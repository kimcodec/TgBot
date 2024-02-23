FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./
COPY messages.txt ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /val

EXPOSE 8080

# Run
CMD ["/val"]