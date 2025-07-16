package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	botToken    = "8047920092:AAGDis_dQ1sjwopmR9MXXawrctPh4fNAZ4w"
	chatID      = "1298975161"
	checkPeriod = 60 * time.Second
)

type ApiResponse struct {
	Results json.RawMessage `json:"results"`
}

var notified = make(map[string]bool)

func sendTelegram(msg string) {
	tgURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", msg)
	data.Set("parse_mode", "Markdown")
	data.Set("disable_web_page_preview", "true")

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", tgURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("❌ Telegram isteği oluşturulamadı: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Telegram isteği başarısız: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Telegram API hatası: %s\nYanıt: %s", resp.Status, string(body))
		return
	}

	log.Printf("✅ Telegram mesajı gönderildi. Yanıt: %s", string(body))
}

func fetchAndProcess() {
	baseURL := "https://www.tesla.com/coinorder/api/v4/inventory-results"

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
			"lng":          28.9533,
			"lat":          41.0145,
			"zip":          "34791",
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

    req.Header.Set("accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`)
    req.Header.Set("accept-encoding", `gzip, deflate, br, zstd`)
    req.Header.Set("accept-language", `tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7`)
    req.Header.Set("cache-control", `no-cache`)
    req.Header.Set("pragma", `no-cache`)
    req.Header.Set("priority", `u=0, i`)
    req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"`)
    req.Header.Set("sec-ch-ua-mobile", `?0`)
    req.Header.Set("sec-ch-ua-platform", `"macOS"`)
    req.Header.Set("sec-fetch-dest", `document`)
    req.Header.Set("sec-fetch-mode", `navigate`)
    req.Header.Set("sec-fetch-site", `none`)
    req.Header.Set("sec-fetch-user", `?1`)
    req.Header.Set("upgrade-insecure-requests", `1`)
    req.Header.Set("user-agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36`)

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

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("❌ JSON parse hatası (root): %v", err)
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

	if err := json.Unmarshal(apiResp.Results, &results); err != nil {
		log.Printf("📋 Uyarı: results parse edilemedi: %v", err)
		return
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
	log.Println("📈 Tesla MYRWD bot başlıyor…")
	ticker := time.NewTicker(checkPeriod)
	defer ticker.Stop()

	for {
		fetchAndProcess()
		<-ticker.C
	}
}
