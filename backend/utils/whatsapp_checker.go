package utils

import (
	"SIAku/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WhatsAppStatus struct {
	Success   bool   `json:"success"`
	Service   string `json:"service"`
	Status    string `json:"status"`
	Connected bool   `json:"connected"`
	BotName   string `json:"bot_name"`
	BotNumber string `json:"bot_number"`
	Port      int    `json:"port"`
}

// CheckWhatsAppService checks if WhatsApp service is running
func CheckWhatsAppService() {
	whatsappURL := config.AppConfig.WhatsAppServiceURL
	if whatsappURL == "" {
		whatsappURL = "http://localhost:3000" // fallback
	}
	
	fmt.Println("\n🔍 Checking WhatsApp Bot Service...")
	
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	
	resp, err := client.Get(whatsappURL)
	if err != nil {
		fmt.Printf("❌ WhatsApp Bot Service: OFFLINE (%s)\n", whatsappURL)
		fmt.Printf("   Error: %v\n", err)
		fmt.Println("   💡 Tip: Start WhatsApp service with 'npm start' in whatsapp folder")
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ WhatsApp Bot Service: ERROR (Status: %d)\n", resp.StatusCode)
		return
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read response: %v\n", err)
		return
	}
	
	var status WhatsAppStatus
	if err := json.Unmarshal(body, &status); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		return
	}
	
	// Print status
	if status.Connected {
		fmt.Printf("✅ WhatsApp Bot Service: ONLINE & CONNECTED\n")
		fmt.Printf("   📱 Status: %s\n", status.Status)
		fmt.Printf("   🤖 Bot Name: %s\n", status.BotName)
		fmt.Printf("   📞 Bot Number: %s\n", status.BotNumber)
	} else {
		fmt.Printf("⚠️  WhatsApp Bot Service: ONLINE but NOT CONNECTED\n")
		fmt.Printf("   📱 Status: %s\n", status.Status)
		fmt.Printf("   🔗 URL: %s\n", whatsappURL)
		fmt.Println("   💡 Scan QR code to connect WhatsApp")
	}
}
