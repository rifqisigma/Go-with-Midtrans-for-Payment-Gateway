package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

func VoidTransaction(orderID string) error {
	url := fmt.Sprintf("https://api.sandbox.midtrans.com/v2/%s/void", orderID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(os.Getenv("MIDTRANS_SERVER_KEY") + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("void gagal: %s", string(body))
	}
	return nil
}
