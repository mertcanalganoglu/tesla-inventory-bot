# Tesla Inventory Telegram Bot 🚗

Tesla envanterinde sadece **Standard Range** olan araçları takip eder, renk, fiyat, VIN ve sipariş linkiyle birlikte Telegram’a bildirir.

---

## 🚀 Kurulum

### 1️⃣ Repo’yu klonla
```
git clone <senin github linkin>
cd tesla-inventory-bot
```

### 2️⃣ Go modüllerini yükle
```
go mod tidy
```

### 3️⃣ `main.go` içinde bot token ve chat ID kontrol et
```go
const (
	botToken = "8047920092:AAGDis_dQ1sjwopmR9MXXawrctPh4fNAZ4w"
	chatID   = "8047920092"
)
```

---

## 🏃 Çalıştır
```
go run main.go
```

veya derleyip binary oluştur:
```
go build -o tesla-bot
./tesla-bot
```

---

## ⏰ Özellikler
✅ Sadece “Standard” geçen araçları bildirir  
✅ Renk, fiyat, VIN ve sipariş linkini gönderir  
✅ 5 dakikada bir kontrol eder  
✅ Tek binary ile çalışır  

---

## 📋 Notlar
- Tesla’nın HTML yapısı değişirse `.Find()` seçicileri güncellemen gerekebilir.
- Botun sana mesaj atabilmesi için önce ona `/start` yazmalısın.
