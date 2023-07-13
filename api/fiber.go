package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type FiberAPI struct {
	url    string
	apiKey string
}

func NewFiberAPI(url string, apiKey string) *FiberAPI {
	url = url + "/api/fiber"
	return &FiberAPI{url, apiKey}
}

type Response[T any] struct {
	Status  string `json:"status"`
	Message T      `json:"message"`
	Error   string `json:"error"`
}

type Quota struct {
	EgressMB         int `json:"egress_mb"`
	MaxEgressMB      int `json:"max_egress_mb"`
	TransactionCount int `json:"transaction_count"`
	ActiveStreams    int `json:"active_streams"`
	MaxActiveStreams int `json:"max_active_streams"`
}

func (f *FiberAPI) prepareRequest(req *http.Request) {
	req.Header.Set("x-api-key", f.apiKey)
	req.Header.Set("Content-Type", "application/json")
}

func (f *FiberAPI) GetQuota() (*Quota, error) {
	req, err := http.NewRequest("GET", f.url+"/quota", nil)
	if err != nil {
		return nil, err
	}

	f.prepareRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := new(Response[*Quota])
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("Error getting quota: %s", response.Error)
	}

	return response.Message, nil
}

type TraceEntry struct {
	Timestamp       uint64 `json:"timestamp"`
	NodeID          string `json:"node_id"`
	Region          string `json:"region"`
	Source          string `json:"source"`
	ObservationType string `json:"observation_type"`
}

func (f *FiberAPI) TraceTransaction(hash, observationType string, private bool) ([]*TraceEntry, error) {
	req, err := http.NewRequest("GET", f.url+"/trace/tx/"+hash, nil)
	if err != nil {
		return nil, err
	}

	if observationType == "" {
		observationType = "all"
	}

	q := req.URL.Query()
	q.Add("observation_type", observationType)
	q.Add("private", strconv.FormatBool(private))
	req.URL.RawQuery = q.Encode()

	f.prepareRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := new(Response[[]*TraceEntry])
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("Error getting quota: %s", response.Error)
	}

	return response.Message, nil
}

func (f *FiberAPI) TraceBlock(hashOrNumber, observationType string) ([]*TraceEntry, error) {
	req, err := http.NewRequest("GET", f.url+"/trace/block/"+hashOrNumber, nil)
	if err != nil {
		return nil, err
	}

	if observationType == "" {
		observationType = "all"
	}

	q := req.URL.Query()
	q.Add("observation_type", observationType)
	req.URL.RawQuery = q.Encode()

	f.prepareRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := new(Response[[]*TraceEntry])
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("Error getting quota: %s", response.Error)
	}

	return response.Message, nil
}
