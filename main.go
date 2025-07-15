package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	apiURL      = "https://www.tesla.com/coinorder/api/v4/inventory-results"
	botToken    = "8047920092:AAGDis_dQ1sjwopmR9MXXawrctPh4fNAZ4w"
	chatID      = "8047920092"
	checkPeriod = 6 * time.Second
)

var seen = make(map[string]bool)

type ApiResponse struct {
	Results []struct {
		VIN           string  `json:"VIN"`
		InventoryID   string  `json:"InventoryID"`
		Price         float64 `json:"Price"`
		TrimName      string  `json:"TrimName"`
		ExteriorColor string  `json:"ExteriorColor"`
	} `json:"results"`
}

func sendTelegram(msg string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	_, err := http.PostForm(apiURL, url.Values{
		"chat_id":    {chatID},
		"text":       {msg},
		"parse_mode": {"Markdown"},
	})
	if err != nil {
		log.Println("Telegram gönderim hatası:", err)
	}
}

func fetchInventory() {
	// JSON gövde
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"model":       "my",
			"condition":   "new",
			"options":     map[string]interface{}{},
			"arrangeby":   "Price",
			"order":       "asc",
			"market":      "TR",
			"language":    "tr",
			"super_region": "north america",
			"lng":         28.9533,
			"lat":         41.0145,
			"zip":         "34096",
			"range":       0,
		},
		"offset":                          0,
		"count":                           24,
		"outsideOffset":                   0,
		"outsideSearch":                   false,
		"isFalconDeliverySelectionEnabled": true,
		"version":                         "v2",
	}

	body, _ := json.Marshal(query)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Request oluşturulamadı:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("API hatası:", err)
		return
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)

	var data ApiResponse
	if err := json.Unmarshal(resBody, &data); err != nil {
		log.Println("JSON parse hatası:", err)
		return
	}

	for _, car := range data.Results {
		if car.TrimName != "MYRWD" {
			continue
		}

		msg := fmt.Sprintf(
			"🚗 *%s*\n💰 *Fiyat:* %.0f ₺\n🎨 *Renk:* %s\n🔢 *VIN:* %s\n\n🔗 [Sipariş Et](https://www.tesla.com/tr_tr/my/order/%s)",
			car.TrimName, car.Price, car.ExteriorColor, car.VIN, car.InventoryID,
		)

		if seen[car.VIN] {
			continue
		}

		log.Println("Yeni MYRWD bulundu:", car.VIN)
		sendTelegram(msg)
		seen[car.VIN] = true
	}
}

func main() {
	log.Println("Tesla API bot başlıyor…")
	ticker := time.NewTicker(checkPeriod)
	defer ticker.Stop()

	for {
		fetchInventory()
		<-ticker.C
	}
}
