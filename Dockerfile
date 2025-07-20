FROM golang:1.21-alpine AS builder

WORKDIR /app

# Go modüllerini kopyala
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Binary'yi build et
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Binary'yi kopyala
COPY --from=builder /app/main .

# Ana uygulamayı çalıştır
CMD ["./main"] 