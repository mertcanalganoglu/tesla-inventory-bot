package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiURL      = "https://www.tesla.com/coinorder/api/v4/inventory-results?query=%7B%22query%22%3A%7B%22model%22%3A%22my%22%2C%22condition%22%3A%22new%22%2C%22options%22%3A%7B%7D%2C%22arrangeby%22%3A%22Price%22%2C%22order%22%3A%22asc%22%2C%22market%22%3A%22TR%22%2C%22language%22%3A%22tr%22%2C%22super_region%22%3A%22north%20america%22%2C%22lng%22%3A28.9533%2C%22lat%22%3A41.0145%2C%22zip%22%3A%2234096%22%2C%22range%22%3A0%7D%2C%22offset%22%3A0%2C%22count%22%3A24%2C%22outsideOffset%22%3A0%2C%22outsideSearch%22%3Afalse%2C%22isFalconDeliverySelectionEnabled%22%3Atrue%2C%22version%22%3A%22v2%22%7D"
	botToken    = "8047920092:AAGDis_dQ1sjwopmR9MXXawrctPh4fNAZ4w"
	chatID      = "8047920092"
	checkPeriod = 60 * time.Second
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
		if strings.TrimSpace(car.TrimName) != "MYRWD" {
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
