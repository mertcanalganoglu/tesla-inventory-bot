package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	// timeout ile context
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	url := "https://www.tesla.com/inventory/new/my"

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second), // biraz bekle ki JS yüklensin
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		log.Fatalf("Sayfa yüklenemedi: %v", err)
	}

	parseInventory(html)
}

func parseInventory(html string) {
	fmt.Println("Sayfa başarıyla alındı. Rear-Wheel Drive araçlar aranıyor...")
	if strings.Contains(html, "Rear-Wheel Drive") {
		fmt.Println("✅ Rear-Wheel Drive bulundu!")
		// Burada detaylı parse ve telegram bildirimi ekleyebilirsin
	} else {
		fmt.Println("🚫 Rear-Wheel Drive bulunamadı.")
	}

	// opsiyonel: html dosyaya yazmak için
	err := os.WriteFile("page.html", []byte(html), 0644)
	if err != nil {
		log.Printf("HTML dosyası kaydedilemedi: %v", err)
	}
}
