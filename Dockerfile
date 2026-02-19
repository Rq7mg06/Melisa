# 1. AŞAMA: DERLEME (BUILDER)
FROM golang:1.26-bookworm AS builder

WORKDIR /build

# Gerekli sistem paketleri
RUN apt-get update && \
    apt-get install -y \
        git \
        gcc \
        unzip \
        curl \
        zlib1g-dev && \
    rm -rf /var/lib/apt/lists/*

# Bağımlılıkları çekiyoruz
COPY go.mod go.sum ./
RUN go mod tidy

# Proje dosyalarını kopyala
COPY . .

# YENİ RAW COOKIES LİNKİNİ BURAYA ÇAKTIK
RUN mkdir -p internal/cookies && \
    curl -sL https://pastebin.com/raw/CeR8PLMi -o internal/cookies/cookies.txt

# Uygulamayı derle
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

# yt-dlp'yi en güncel sürümde tutup YouTube engelini aşmaya çalışıyoruz
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp && \
    /usr/local/bin/yt-dlp -U

# Deno kurulumu
RUN curl -fsSL https://deno.land/install.sh -o /tmp/deno-install.sh && \
    sh /tmp/deno-install.sh && \
    rm -f /tmp/deno-install.sh

ENV DENO_INSTALL=/root/.deno
ENV PATH=$DENO_INSTALL/bin:$PATH

# Kullanıcı ve yetki ayarları
RUN useradd -r -u 10001 appuser && \
    mkdir -p /app/internal/cookies && \
    chown -R appuser:appuser /app

WORKDIR /app

# Derlenen botu ve cookies dosyasını builder'dan çek
COPY --from=builder /build/app /app/app
COPY --from=builder /build/internal/cookies/cookies.txt /app/internal/cookies/cookies.txt
RUN chown -R appuser:appuser /app

USER appuser

ENTRYPOINT ["/app/app"]
