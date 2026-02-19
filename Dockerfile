# 1. AŞAMA: DERLEME (BUILDER)
# Go sürümünü 1.26 yaparak o meşhur "1.25.7 hatasını" kökten çözüyoruz.
FROM golang:1.26-bookworm AS builder

WORKDIR /build

# Müzik botu için gerekli olan derleme araçlarını kuruyoruz.
RUN apt-get update && \
    apt-get install -y \
        git \
        gcc \
        unzip \
        curl \
        zlib1g-dev && \
    rm -rf /var/lib/apt/lists/*

# Bağımlılıkları (go.mod ve go.sum) kopyalayıp güncelliyoruz.
COPY go.mod go.sum ./
RUN go mod tidy

# Tüm proje dosyalarını içeri alıp derlemeyi başlatıyoruz.
COPY . .

RUN chmod +x install.sh && \
    ./install.sh -n --quiet --skip-summary && \
    CGO_ENABLED=1 go build -v -trimpath -ldflags="-w -s" -o app ./cmd/app/


# 2. AŞAMA: ÇALIŞTIRMA (FINAL IMAGE)
FROM debian:bookworm-slim

# Ses ve video işleme için ffmpeg ve yt-dlp gibi kritik araçlar burada kurulur.
RUN apt-get update && \
    apt-get install -y \
        ffmpeg \
        curl \
        unzip \
        zlib1g && \
    rm -rf /var/lib/apt/lists/*

# SSL sertifikalarını güvenli bağlantı için builder'dan çekiyoruz.
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

# YouTube videolarını indirmek için yt-dlp ve müzik motoru için Deno kurulumu.
RUN curl -fL \
      https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux \
      -o /usr/local/bin/yt-dlp && \
    chmod 0755 /usr/local/bin/yt-dlp && \
    curl -fsSL https://deno.land/install.sh -o /tmp/deno-install.sh && \
    sh /tmp/deno-install.sh && \
    rm -f /tmp/deno-install.sh

ENV DENO_INSTALL=/root/.deno
ENV PATH=$DENO_INSTALL/bin:$PATH

# Güvenlik için kullanıcı ayarları.
RUN useradd -r -u 10001 appuser && \
    mkdir -p /app && \
    chown -R appuser:appuser /app

WORKDIR /app

# Derlenen müzik botu dosyasını kopyalıyoruz.
COPY --from=builder /build/app /app/app
RUN chown appuser:appuser /app/app

USER appuser

# Botu ateşleyen komut.
ENTRYPOINT ["/app/app"]
