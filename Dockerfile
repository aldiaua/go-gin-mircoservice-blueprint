# --- Stage 1: Build Binary ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies yang dibutuhkan untuk build (jika diperlukan)
RUN apk add --no-cache git

# Salin modul go
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode sumber
COPY . .

# Build aplikasi Go (sesuaikan path entrypoint di ./cmd/users jika berbeda)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/service ./cmd/users

# --- Stage 2: Run Image ---
FROM alpine:latest

WORKDIR /app

# Salin binary hasil build dari stage sebelumnya
COPY --from=builder /app/bin/service /app/service

# Salin direktori migrations jika dibutuhkan oleh aplikasi/runtime
COPY --from=builder /app/migrations /app/migrations

# Ekspose port aplikasi (sesuaikan jika port service Anda berbeda)
EXPOSE 8080

# Jalankan service
CMD ["/app/service"]