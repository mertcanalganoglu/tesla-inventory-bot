package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	botToken    = "8047920092:AAGDis_dQ1sjwopmR9MXXawrctPh4fNAZ4w"
	chatID      = "8047920092"
	checkPeriod = 60 * time.Second
)

var seen = make(map[string]bool)

type ApiResponse struct {
	Results []struct {
		VIN         string   `json:"VIN"`
		Price       float64  `json:"Price"`
		TrimName    string   `json:"TrimName"`
		TRIM        []string `json:"TRIM"`
		InventoryID string   `json:"InventoryID"`
	} `json:"results"`
}

func buildURL() string {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"model":        "my",
			"condition":    "new",
			"options":      map[string]interface{}{},
			"arrangeby":    "Price",
			"order":        "asc",
			"market":       "DE",
			"language":     "de",
			"super_region": "north america",
		},
		"offset":                          0,
		"count":                           24,
		"outsideOffset":                   0,
		"outsideSearch":                   false,
		"isFalconDeliverySelectionEnabled": true,
		"version":                         "v2",
	}

	j, err := json.Marshal(query)
	if err != nil {
		log.Fatal(err)
	}

	escaped := url.QueryEscape(string(j))
	return fmt.Sprintf("https://www.tesla.com/inventory/api/v4/inventory-results?query=%s", escaped)
}

func sendTelegram(msg string) {
	tgURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	_, err := http.PostForm(tgURL, url.Values{
		"chat_id":    {chatID},
		"text":       {msg},
		"parse_mode": {"Markdown"},
	})
	if err != nil {
		log.Println("Telegram gönderim hatası:", err)
	}
}

func fetchInventory(apiURL string) {
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Println("API hatası:", err)
		return
	}
	defer resp.Body.Close()

	var data ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("JSON parse hatası:", err)
		return
	}

	for _, car := range data.Results {
		foundMYRWD := false
		for _, t := range car.TRIM {
			if t == "MYRWD" {
				foundMYRWD = true
				break
			}
		}
		if !foundMYRWD {
			continue
		}

		if seen[car.VIN] {
			continue
		}

		msg := fmt.Sprintf(
			"🚗 *%s*\n💰 *Fiyat:* %.0f EUR\n🔢 *VIN:* %s\n\n🔗 [Sipariş Et](https://www.tesla.com/de_DE/my/order/%s)",
			car.TrimName, car.Price, car.VIN, car.InventoryID,
		)
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
		apiURL := buildURL()
		fetchInventory(apiURL)
		<-ticker.C
	}
}
