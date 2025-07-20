package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Telegram bot token ve chat ID
var (
	botToken = "YOUR_BOT_TOKEN_HERE" // Buraya bot token'ınızı yazın
	chatID   = int64(0)              // Buraya chat ID'nizi yazın
)

// Session cookie'leri
var sessionCookies []*http.Cookie

// Bildirilen araçları takip et
var notified = make(map[string]bool)

// ApiResponse yapısı
type ApiResponse struct {
	Results struct {
		Exact              []json.RawMessage `json:"exact"`
		Approximate        []json.RawMessage `json:"approximate"`
		ApproximateOutside []json.RawMessage `json:"approximateOutside"`
	} `json:"results"`
	TotalMatchesFound int `json:"total_matches_found"`
}

// User-Agent rotasyonu
var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
}

func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func sendTelegram(msg string) {
	if botToken == "YOUR_BOT_TOKEN_HERE" || chatID == 0 {
		log.Printf("⚠️ Telegram bot token veya chat ID ayarlanmamış")
		return
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("❌ Telegram bot oluşturulamadı: %v", err)
		return
	}

	message := tgbotapi.NewMessage(chatID, msg)
	message.ParseMode = "Markdown"
	message.DisableWebPagePreview = true

	_, err = bot.Send(message)
	if err != nil {
		log.Printf("❌ Telegram mesajı gönderilemedi: %v", err)
	} else {
		log.Printf("✅ Telegram mesajı gönderildi")
	}
}

// Health check endpoint'i
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"OK","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

func getSessionCookies() error {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", "https://www.tesla.com/tr_TR/inventory", nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", getRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sessionCookies = resp.Cookies()
	log.Printf("✅ Session cookie'leri alındı: %d adet", len(sessionCookies))
	for i, cookie := range sessionCookies {
		log.Printf("🍪 Cookie %d: %s", i+1, cookie.Name+"="+cookie.Value)
	}

	return nil
}

func fetchAndProcess() {
	// Session cookie'leri yoksa al
	if len(sessionCookies) == 0 {
		if err := getSessionCookies(); err != nil {
			log.Printf("❌ Session cookie'leri alınamadı: %v", err)
			return
		}
	}

	// Gerçek browser isteğindeki URL'yi kullan
	baseURL := "https://www.tesla.com/coinorder/api/v4/inventory-results"

	// Gerçek browser isteğindeki query yapısını kullan
	queryPayload := map[string]interface{}{
		"query": map[string]interface{}{
			"model":        "my",
			"condition":    "new",
			"options":      map[string]interface{}{},
			"arrangeby":    "Price",
			"order":        "asc",
			"market":       "TR",
			"language":     "tr",
			"super_region": "north america",
			"lng":          "",
			"lat":          "",
			"zip":          "",
			"range":        0,
		},
		"offset":                           0,
		"count":                            24,
		"outsideOffset":                    0,
		"outsideSearch":                    false,
		"isFalconDeliverySelectionEnabled": true,
		"version":                          "v2",
	}

	queryJSON, err := json.Marshal(queryPayload)
	if err != nil {
		log.Printf("❌ Query JSON oluşturulamadı: %v", err)
		return
	}

	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("query", string(queryJSON))
	u.RawQuery = q.Encode()

	log.Printf("🌐 İstek atılacak URL: %s", u.String())

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Printf("❌ GET isteği oluşturulamadı: %v", err)
		return
	}

	// Session cookie'lerini ekle
	for _, cookie := range sessionCookies {
		req.AddCookie(cookie)
	}

	// Gerçek browser headers'ı kullan
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="137", "Chromium";v="137", "Not/A)Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	log.Println("🚀 GET isteği gönderiliyor...")
	httpResp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ HTTP isteği başarısız oldu: %v", err)
		return
	}
	defer httpResp.Body.Close()

	var reader io.Reader
	switch httpResp.Header.Get("Content-Encoding") {
	case "gzip":
		gzipReader, err := gzip.NewReader(httpResp.Body)
		if err != nil {
			log.Printf("❌ gzip açma hatası: %v", err)
			return
		}
		defer gzipReader.Close()
		reader = gzipReader
	default:
		reader = httpResp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("❌ Yanıt okuma hatası: %v", err)
		return
	}

	log.Printf("🔷 Sunucudan gelen yanıt:\n%s", string(body))

	// Access Denied kontrolü
	if strings.Contains(string(body), "Access Denied") {
		log.Printf("❌ Access Denied - Cookie'ler yenileniyor...")
		sessionCookies = nil         // Cookie'leri temizle
		time.Sleep(10 * time.Second) // Daha uzun bekleme
		return
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("❌ JSON parse hatası (root): %v", err)
		return
	}

	log.Printf("🔍 Toplam eşleşme: %d", apiResp.TotalMatchesFound)

	// Tüm araçları birleştir
	var allVehicles []json.RawMessage
	allVehicles = append(allVehicles, apiResp.Results.Exact...)
	allVehicles = append(allVehicles, apiResp.Results.Approximate...)
	allVehicles = append(allVehicles, apiResp.Results.ApproximateOutside...)

	if len(allVehicles) == 0 {
		log.Printf("ℹ️ Envanterde araç bulunamadı")

		// Araba yoksa da Telegram'a bildirim gönder
		msg := "🔍 *Tesla Envanter Kontrolü*\n\n❌ Şu anda envanterde araç bulunamadı.\n\n⏰ Kontrol zamanı: " + time.Now().Format("15:04:05")
		sendTelegram(msg)
		return
	}

	var results []struct {
		VIN            string   `json:"VIN"`
		InventoryPrice float64  `json:"InventoryPrice"`
		TrimName       string   `json:"TrimName"`
		TRIM           []string `json:"TRIM"`
		InventoryID    string   `json:"InventoryID"`
		PAINT          []string `json:"PAINT"`
		INTERIOR       []string `json:"INTERIOR"`
	}

	// Her araç için parse et
	for _, vehicle := range allVehicles {
		var car struct {
			VIN            string   `json:"VIN"`
			InventoryPrice float64  `json:"InventoryPrice"`
			TrimName       string   `json:"TrimName"`
			TRIM           []string `json:"TRIM"`
			InventoryID    string   `json:"InventoryID"`
			PAINT          []string `json:"PAINT"`
			INTERIOR       []string `json:"INTERIOR"`
		}

		if err := json.Unmarshal(vehicle, &car); err != nil {
			log.Printf("📋 Araç parse edilemedi: %v", err)
			continue
		}

		results = append(results, car)
	}

	for _, car := range results {
		foundMYRWD := false
		for _, t := range car.TRIM {
			if strings.EqualFold(t, "MYRWD") {
				foundMYRWD = true
				break
			}
		}
		if !foundMYRWD || notified[car.VIN] {
			continue
		}

		paint := "Bilinmiyor"
		if len(car.PAINT) > 0 {
			paint = car.PAINT[0]
		}
		interior := "Bilinmiyor"
		if len(car.INTERIOR) > 0 {
			interior = car.INTERIOR[0]
		}

		// Siyah dışındaki renkleri filtrele
		paintLower := strings.ToLower(paint)
		if strings.Contains(paintLower, "black") || strings.Contains(paintLower, "siyah") {
			log.Printf("⚫ Siyah araç atlandı: %s (%s)", car.VIN, paint)
			continue
		}

		id := car.InventoryID
		if id == "" {
			id = car.VIN
		}

		orderLink := fmt.Sprintf(
			"https://www.tesla.com/tr_TR/my/order/%s?titleStatus=new&redirect=no#payment",
			id,
		)

		msg := fmt.Sprintf(
			`🟢 Araç Eklendi: Yeni Model Y (_Model Y Arkadan Çekiş_)

🚘 *Dış Renk:* %s
🎨 *İç Renk:* %s
🔢 *VIN:* %s
💰 *Fiyat:* %.0f TL

🔗 [Sipariş Linki](%s)`,
			escapeMarkdown(paint),
			escapeMarkdown(interior),
			escapeMarkdown(car.VIN),
			car.InventoryPrice,
			orderLink,
		)

		log.Println("✅ MYRWD bulundu ve bildirildi:", car.VIN)
		sendTelegram(msg)
		notified[car.VIN] = true
	}
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
	)
	return replacer.Replace(text)
}

func main() {
	// Random seed
	rand.Seed(time.Now().UnixNano())

	log.Println("📈 Tesla MYRWD bot başlıyor…")
	log.Println("⚙️ Zamanlama: 18:00-19:00 (UTC+3) arası 5 saniyede bir, diğer zamanlarda saatte 1 kontrol")

	for {
		fetchAndProcess()

		// Europe/Istanbul time zone'u yoksa UTC+3 offset'iyle manuel oluştur
		loc, err := time.LoadLocation("Europe/Istanbul")
		if err != nil {
			loc = time.FixedZone("UTC+3", 3*60*60)
		}
		now := time.Now().In(loc)
		hour := now.Hour()

		if hour == 18 {
			time.Sleep(5 * time.Second)
		} else {
			nextHour := now.Truncate(time.Hour).Add(time.Hour)
			dur := time.Until(nextHour)
			log.Printf("⏳ Sonraki kontrol %s sonra (saat başı)", dur.Round(time.Second))
			time.Sleep(dur)
		}
	}
}
