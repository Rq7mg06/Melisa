# 1. AŞAMA: DERLEME (BUILDER)
FROM golang:1.26-bookworm AS builder

WORKDIR /build

RUN apt-get update && \
    apt-get install -y \
        git \
        gcc \
        unzip \
        curl \
        zlib1g-dev && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Çerezleri klasöre çekiyoruz
RUN mkdir -p internal/cookies && \
    curl -sL https://pastebin.com/raw/b9VkXvX4 -o internal/cookies/cookies.txt

RUN chmod +x install.sh && \
    ./install.sh -n --quiet && \
    CGO_ENABLED=1 go build -v -trimpath -ldflags="-w -s" -o app ./cmd/app/


# 2. AŞAMA: ÇALIŞTIRMA (FINAL IMAGE)
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y \
        ffmpeg \
        curl \
        unzip \
        zlib1g \
        ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# YT-DLP'Yİ GÜNCELLEYİP YOUTUBE'U KANDIRMAYA ÇALIŞIYORUZ
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp && \
    /usr/local/bin/yt-dlp -U

# Deno kurulumu
RUN curl -fsSL https://deno.land/install.sh -o /tmp/deno-install.sh && \
    sh /tmp/deno-install.sh && \
    rm -f /tmp/deno-install.sh

ENV DENO_INSTALL=/root/.deno
ENV PATH=$DENO_INSTALL/bin:$PATH

RUN useradd -r -u 10001 appuser && \
    mkdir -p /app/internal/cookies && \
    chown -R appuser:appuser /app

WORKDIR /app

COPY --from=builder /build/app /app/app
COPY --from=builder /build/internal/cookies/cookies.txt /app/internal/cookies/cookies.txt
RUN chown -R appuser:appuser /app

USER appuser

# KRİTİK: YouTube'un IPv6 engeline takılmamak için botu IPv4 zorlamaya çalışabiliriz
# Eğer botun içinde ayar varsa Heroku'dan FORCE_IPV4=true yapmayı unutma.

ENTRYPOINT ["/app/app"]
