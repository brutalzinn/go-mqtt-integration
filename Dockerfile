FROM golang:latest

RUN apt-get update && \
    apt-get install -y \
    youtube-dl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /main
CMD ["/main"]
EXPOSE 80