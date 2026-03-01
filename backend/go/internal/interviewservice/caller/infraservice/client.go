package infraservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wisaitas/cloud-forest-assignment/pkg/httpx"
)

type InfraServiceCaller interface {
	ListSKUs(ctx context.Context) (*SKUsResponse, error)
	Provision(ctx context.Context, sku string) (*ProvisionResponse, error)
	Power(ctx context.Context, resourceID, action string) (*PowerResponse, error)
	IsValidSKU(ctx context.Context, sku string) (bool, error)
}

type infraServiceCaller struct {
	baseURL string
}

func NewClient(baseURL string) InfraServiceCaller {
	return &infraServiceCaller{baseURL: baseURL}
}

func (c *infraServiceCaller) ListSKUs(ctx context.Context) (*SKUsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/skus", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpx.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list skus: status %d", resp.StatusCode)
	}
	var out SKUsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *infraServiceCaller) Provision(ctx context.Context, sku string) (*ProvisionResponse, error) {
	body, err := json.Marshal(ProvisionRequest{SKU: sku})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/resources", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpx.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provision: status %d", resp.StatusCode)
	}
	var out ProvisionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *infraServiceCaller) Power(ctx context.Context, resourceID, action string) (*PowerResponse, error) {
	body, _ := json.Marshal(PowerRequest{Action: action})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/resources/"+resourceID+"/power", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpx.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("power: status %d", resp.StatusCode)
	}
	var out PowerResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *infraServiceCaller) IsValidSKU(ctx context.Context, sku string) (bool, error) {
	list, err := c.ListSKUs(ctx)
	if err != nil {
		return false, err
	}
	for _, s := range list.SKUs {
		if s.SKU == sku {
			return true, nil
		}
	}
	return false, nil
}
