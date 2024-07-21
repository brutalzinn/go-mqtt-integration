FROM golang:latest

WORKDIR /app

RUN echo 'deb https://deb.debian.org/debian stable non-free contrib' >> /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y gcc \
    build-essential \
    ffmpeg \
    python3-pip \
    libtool \
    git \
    libasound2-dev \
    yt-dlp \
    make


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /main .

EXPOSE 3000

ENTRYPOINT ["/main"]
