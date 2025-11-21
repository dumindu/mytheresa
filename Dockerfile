FROM golang:1.25-alpine
WORKDIR /app

RUN apk --update-cache add gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/server ./cmd/server \
    && go build -ldflags '-w -s' -a -o ./bin/seed ./cmd/seed

EXPOSE 8080
CMD ["/app/bin/server"]
