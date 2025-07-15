# Tesla Inventory Bot 🚗 (chromedp)

Tesla envanterinde Rear-Wheel Drive araçları kontrol eder.  
Headless Chrome kullanarak Cloudflare & JS engellerini aşar.

---

## 🚀 Gereksinimler
✅ Go ≥ 1.20  
✅ Chrome veya Chromium yüklü

---

## 🔧 Kurulum

### 1️⃣ Repo'yu klonla
```
git clone <senin-github-repon>
cd tesla-inventory-bot-chromedp
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
✅ Headless tarayıcı ile sayfayı yükler  
✅ Cloudflare & bot korumalarına takılmaz  
✅ Rear-Wheel Drive geçen içerikleri arar  
✅ HTML'i `page.html` olarak kaydeder (isteğe bağlı)

---

## 🔗 Notlar
- Daha detaylı parse ve Telegram bildirimi için `parseInventory()` fonksiyonunu genişletebilirsin.
- 45s timeout ile çalışır.
