FROM golang:1.25.1-alpine
WORKDIR /app
RUN apk add --no-cache bash ffmpeg curl python3
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
RUN go build -o app .
CMD ["./app"]
