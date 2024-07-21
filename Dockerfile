FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o go-schedule-manager .

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

ENTRYPOINT ["./main"]

EXPOSE 8000