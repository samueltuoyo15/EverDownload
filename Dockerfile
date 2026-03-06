FROM golang:1.24-alpine
WORKDIR /app

RUN apk add --no-cache bash ffmpeg curl python3 py3-pip

# Install yt-dlp binary (pinned to a stable release)
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/server/

EXPOSE 5000
CMD ["./app"]
