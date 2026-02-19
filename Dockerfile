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

# Bağımlılıkları çek (Burada artık Go sürüm hatası almayacaksın)
COPY go.mod go.sum ./
RUN go mod tidy

# Tüm proje dosyalarını kopyala
COPY . .

# Hatalı olan --skip-summary komutu kaldırıldı
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
        zlib1g && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# Müzik motoru bileşenleri (yt-dlp ve Deno)
RUN curl -fL \
      https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux \
      -o /usr/local/bin/yt-dlp && \
    chmod 0755 /usr/local/bin/yt-dlp && \
    curl -fsSL https://deno.land/install.sh -o /tmp/deno-install.sh && \
    sh /tmp/deno-install.sh && \
    rm -f /tmp/deno-install.sh

ENV DENO_INSTALL=/root/.deno
ENV PATH=$DENO_INSTALL/bin:$PATH

# Kullanıcı ve dizin ayarları
RUN useradd -r -u 10001 appuser && \
    mkdir -p /app && \
    chown -R appuser:appuser /app

WORKDIR /app

# Uygulamayı builder'dan çek
COPY --from=builder /build/app /app/app
RUN chown appuser:appuser /app/app

USER appuser

ENTRYPOINT ["/app/app"]
