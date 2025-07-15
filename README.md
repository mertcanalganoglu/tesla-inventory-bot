# Tesla Inventory API Bot 🚗 (MYRWD)

Tesla'nın resmi API'sini kullanarak sadece `MYRWD` trimindeki araçları kontrol eder ve Telegram'a bildirir.

---

## 🚀 Gereksinimler
✅ Go ≥ 1.20

---

## 🔧 Kurulum

### 1️⃣ Repo'yu klonla
```
git clone <senin-github-repon>
cd tesla-api-bot-myrwd
```

### 2️⃣ Modülleri yükle
```
go mod tidy
```

---

## 🏃 Çalıştır
```
go run main.go
```

veya binary yap:
```
go build -o tesla-bot
./tesla-bot
```

---

## 📋 Özellikler
✅ Resmi API kullanır  
✅ Sadece `MYRWD` trim olanları filtreler  
✅ Fiyat, renk, VIN, sipariş linki gönderir  
✅ 60 saniyede bir kontrol eder  
✅ Cloudflare & bot engeli yok  
✅ Hızlı & stabil

---

## 🔗 Notlar
- Telegram bot token ve chat ID'yi kodda değiştirmeyi unutma.
