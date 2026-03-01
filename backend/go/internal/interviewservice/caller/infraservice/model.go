package infraservice

type SKUItem struct {
	ID   string `json:"id"`
	SKU  string `json:"sku"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type SKUsResponse struct {
	SKUs []SKUItem `json:"skus"`
}

type ProvisionRequest struct {
	SKU string `json:"sku"`
}

type ProvisionResponse struct {
	ID        string `json:"id"`
	SKU       string `json:"sku"`
	Status    string `json:"status"`
	IP        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}

type PowerRequest struct {
	Action string `json:"action"`
}

type PowerResponse struct {
	Status string `json:"status"`
	State  string `json:"state"`
}
