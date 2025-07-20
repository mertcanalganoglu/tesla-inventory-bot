# Tesla Inventory Bot

Tesla Model Y (MYRWD) envanterini takip eden ve Telegram üzerinden bildirim gönderen bot.

## Özellikler

- ✅ Tesla API'sini kullanarak envanter kontrolü
- ✅ Session cookie yönetimi
- ✅ Bot koruması bypass
- ✅ Telegram bildirimleri
- ✅ Siyah dışındaki renkleri filtreleme
- ✅ 5 saniyede bir kontrol
- ✅ Render.com deployment hazır

## Render.com Deployment

### 1. Repository'yi Render.com'a bağla

1. Render.com'da yeni bir "Web Service" oluştur
2. GitHub repository'yi bağla
3. Build Command: `go build -o main .`
4. Start Command: `./main`

### 2. Environment Variables

Render.com dashboard'unda şu environment variable'ları ekle:

```
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
```

### 3. Bot Token Alma

1. Telegram'da @BotFather'a git
2. `/newbot` komutunu kullan
3. Bot adı ve username ver
4. Bot token'ını al

### 4. Chat ID Alma

1. Bot'u Telegram'da başlat
2. Bir mesaj gönder
3. Bu URL'yi ziyaret et: `https://api.telegram.org/bot<BOT_TOKEN>/getUpdates`
4. Chat ID'yi bul

## Yerel Çalıştırma

```bash
# Dependencies'i yükle
go mod tidy

# Çalıştır
go run main.go
```

## Health Check

Bot `/health` endpoint'i ile sağlık kontrolü yapar:

```bash
curl https://your-app.onrender.com/health
```

## Loglar

Bot detaylı loglar verir:
- Session cookie'leri
- API istekleri
- Bulunan araçlar
- Telegram bildirimleri

## Güvenlik

- Bot token'ları environment variable olarak saklanır
- Session cookie'leri otomatik yenilenir
- User-Agent rotasyonu
- Realistic browser headers
