# Tesla Inventory API Bot 🚗 (MYRWD, URL Fix)

Tesla'nın resmi API'sini doğru encode edilmiş URL ile kullanır, yalnızca `MYRWD` olanları filtreler ve Telegram'a bildirir.

---

## 🚀 Gereksinimler
✅ Go ≥ 1.20

---

## 🔧 Kurulum

### 1️⃣ Repo'yu klonla
```
git clone <senin-github-repon>
cd tesla-api-bot-fixed
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
✅ URL parametresini programatik olarak ve doğru encode eder  
✅ Sadece `MYRWD` olanları filtreler  
✅ Fiyat, VIN, sipariş linki gönderir  
✅ 60 saniyede bir kontrol eder  
✅ Telegram'a bildirir

---

## 🔗 Notlar
- Telegram bot token ve chat ID'yi kodda değiştirmeyi unutma.
