package services

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func SendNotification(message string) {
	target := os.Getenv("NOTIFICATION_TARGET")

	if target == "" {
		fmt.Println("Error: NOTIFICATION_TARGET is missing in .env")
		return
	}

	fmt.Printf("Sending alert to configured target...\n")

	url := fmt.Sprintf("https://ntfy.sh/%s", target)
	http.Post(url, "text/plain", strings.NewReader(message))
}
