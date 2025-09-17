package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (a *api) fetchPrice(productID int64) (int64, error) {
	url, err := url.JoinPath(a.productsURL, "products", strconv.FormatInt(productID, 10))
	if err != nil {
		return 0, fmt.Errorf("join path: %w", err)
	}
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return 0, fmt.Errorf("program not found (404) url: %s", url)
	}
	if resp.StatusCode >= 500 {
		return 0, fmt.Errorf("products-service error %d", resp.StatusCode)
	}
	var body priceLookup
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("decode and encoder error")
	}

	return body.Price, nil
}

func (a *api) reserve(productID, qty int64) (bool, error) {
	b, _ := json.Marshal(
		reserveReq{Qty: qty},
	)

	url, err := url.JoinPath(a.productsURL, "products", strconv.FormatInt(productID, 10), "reserve")
	if err != nil {
		return false, fmt.Errorf("join path: %w", err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return false, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return true, nil
	case http.StatusConflict:
		// insufficient stock
		return false, nil
	case http.StatusNotFound:
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return false, fmt.Errorf("product %d not found: %s", productID, string(b))
	default:
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return false, fmt.Errorf("reserve failed: status=%d url=%s body=%q", resp.StatusCode, url, string(b))
	}
}
