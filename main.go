package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	teslaURL    = "https://www.tesla.com/inventory/new/my"
	botToken    = "8047920092:AAGDis_dQ1sjwopmR9MXXawrctPh4fNAZ4w"
	chatID      = "8047920092"
	checkPeriod = 5 * time.Minute
)

var seen = make(map[string]bool)

func sendTelegram(msg string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id":    {chatID},
		"text":       {msg},
		"parse_mode": {"Markdown"},
	})
	if err != nil {
		log.Println("Telegram gönderim hatası:", err)
		return
	}
	defer resp.Body.Close()
}

func fetchInventory() ([]string, error) {
	resp, err := http.Get(teslaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var vehicles []string
	doc.Find("div[data-test='vehicleCard']").Each(func(i int, s *goquery.Selection) {
		var model, price, color, vin, orderLink string
		isStandard := false

		if title := s.Find("h2").Text(); title != "" {
			model = strings.TrimSpace(title)
			if strings.Contains(strings.ToLower(model), "standard") {
				isStandard = true
			}
		}

		s.Find(".vehicle-attribute").Each(func(j int, attr *goquery.Selection) {
			txt := strings.ToLower(attr.Text())
			if strings.Contains(txt, "standard") {
				isStandard = true
			}
			if strings.HasPrefix(strings.ToLower(txt), "vin") {
				vin = strings.TrimSpace(attr.Text())
			}
		})

		if !isStandard {
			return
		}

		if p := s.Find(".vehicle-price").Text(); p != "" {
			price = strings.TrimSpace(p)
		}

		if c := s.Find(".color-name").Text(); c != "" {
			color = strings.TrimSpace(c)
		}

		if link, exists := s.Find("a[data-test='vehicleCardCTA']").Attr("href"); exists {
			orderLink = fmt.Sprintf("https://www.tesla.com%s", link)
		}

		if vin == "" {
			vinText := s.Find("div:contains('VIN')").Text()
			if vinText != "" {
				vin = strings.TrimSpace(vinText)
			}
		}

		message := fmt.Sprintf(
			"🚗 *%s*\n💰 *Fiyat:* %s\n🎨 *Renk:* %s\n🔢 *VIN:* %s\n\n🔗 [Sipariş Et](%s)",
			model, price, color, vin, orderLink,
		)

		if message != "" && !seen[message] {
			vehicles = append(vehicles, message)
			seen[message] = true
		}
	})
	return vehicles, nil
}

func check() {
	vehicles, err := fetchInventory()
	if err != nil {
		log.Println("Envanter kontrol hatası:", err)
		return
	}

	for _, v := range vehicles {
		log.Println("Yeni *Standard Range* araç bulundu:")
		log.Println(v)
		sendTelegram(v)
	}
	log.Printf("Kontrol tamamlandı. %d araç bildirildi.\n", len(vehicles))
}

func main() {
	log.Println("Tesla *Standard Range* envanter botu başlıyor…")
	check()

	ticker := time.NewTicker(checkPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			check()
		}
	}
}
