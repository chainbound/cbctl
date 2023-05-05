package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type FiberAPI struct {
	url    string
	apiKey string
}

func NewFiberAPI(url string, apiKey string) *FiberAPI {
	url = url + "/api/fiber"
	return &FiberAPI{url, apiKey}
}

type Quota struct {
	EgressMB         int `json:"egress_mb"`
	MaxEgressMB      int `json:"max_egress_mb"`
	TransactionCount int `json:"transaction_count"`
	ActiveStreams    int `json:"active_streams"`
	MaxActiveStreams int `json:"max_active_streams"`
}

func (f *FiberAPI) GetQuota() (*Quota, error) {
	req, err := http.NewRequest("GET", f.url+"/quota", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", f.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	quota := new(Quota)
	if err := json.Unmarshal(body, quota); err != nil {
		return nil, err
	}

	return quota, nil
}
