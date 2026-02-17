package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	port        = "8081"
	failureRate = 10
	slowRate    = 20
)

type SKU struct {
	ID           string  `json:"id"`
	SKU          string  `json:"sku"`
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	CPU          int     `json:"cpu"`
	RAM          int     `json:"ram"`
	Disk         int     `json:"disk"`
	PriceHourly  float64 `json:"price_hourly"`
	PriceMonthly float64 `json:"price_monthly"`
}

var availableSKUs = []SKU{
	{SKU: "C1-R1GB-D20GB", Type: "virtual-machine", Name: "Micro", CPU: 1, RAM: 1, Disk: 20, PriceHourly: 0.27, PriceMonthly: 180},
	{SKU: "C2-R4GB-D80GB", Type: "virtual-machine", Name: "Standard", CPU: 2, RAM: 4, Disk: 80, PriceHourly: 1.1, PriceMonthly: 750},
	{SKU: "C4-R8GB-D160GB", Type: "virtual-machine", Name: "Performance", CPU: 4, RAM: 8, Disk: 160, PriceHourly: 2.2, PriceMonthly: 1500},
	{SKU: "C8-R32GB-D320GB", Type: "virtual-machine", Name: "Pro Max", CPU: 8, RAM: 32, Disk: 320, PriceHourly: 5.2, PriceMonthly: 3500},
	{SKU: "C8-R16GB-D512GB", Type: "dedicated", Name: "Metal Alpha", CPU: 8, RAM: 16, Disk: 512, PriceHourly: 18, PriceMonthly: 12000},
	{SKU: "C16-R64GB-D1024GB", Type: "dedicated", Name: "Metal Beta", CPU: 16, RAM: 64, Disk: 1024, PriceHourly: 42, PriceMonthly: 28000},
	{SKU: "C32-R128GB-D2048GB", Type: "dedicated", Name: "Metal Gamma", CPU: 32, RAM: 128, Disk: 2048, PriceHourly: 90, PriceMonthly: 60000},
	{SKU: "C64-R256GB-D4096GB", Type: "dedicated", Name: "Metal Omega", CPU: 64, RAM: 256, Disk: 4096, PriceHourly: 180, PriceMonthly: 120000},
}

var (
	resources = make(map[string]*Resource)
	mu        sync.RWMutex
)

type Resource struct {
	ID        string `json:"id"`
	SKU       string `json:"sku"`
	Status    string `json:"status"`
	IP        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}

type AvailabilityRequest struct {
	SKU string `json:"sku"`
}

type ProvisionRequest struct {
	SKU string `json:"sku"`
}

type PowerRequest struct {
	Action string `json:"action"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/v1/skus", handleSKUs)
	http.HandleFunc("/v1/availability", handleAvailability)
	http.HandleFunc("/v1/resources", handleResources)
	http.HandleFunc("/v1/resources/", handleResourceAction)

	fmt.Printf("🔌 Infra Service running on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleSKUs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"skus": availableSKUs})
}

func handleAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	exists := false
	for _, s := range availableSKUs {
		if s.SKU == req.SKU {
			exists = true
			break
		}
	}

	available := exists
	if req.SKU == "OUT-OF-STOCK" {
		available = false
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"sku":       req.SKU,
		"available": available,
	})
}

func handleResources(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mu.RLock()
		defer mu.RUnlock()
		list := make([]*Resource, 0, len(resources))
		for _, res := range resources {
			list = append(list, res)
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{"resources": list})
		return
	}

	if r.Method == http.MethodPost {
		simulateChaos(w)

		var req ProvisionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		validSKU := false
		for _, s := range availableSKUs {
			if s.SKU == req.SKU {
				validSKU = true
				break
			}
		}
		if !validSKU && req.SKU != "test" {
			http.Error(w, `{"error": "Invalid SKU"}`, http.StatusBadRequest)
			return
		}

		time.Sleep(time.Duration(500+rand.Intn(1000)) * time.Millisecond)

		resourceID := fmt.Sprintf("i-%d", rand.Int63())
		res := &Resource{
			ID:        resourceID,
			SKU:       req.SKU,
			Status:    "running",
			IP:        fmt.Sprintf("10.0.%d.%d", rand.Intn(255), rand.Intn(255)),
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		mu.Lock()
		resources[resourceID] = res
		mu.Unlock()

		respondJSON(w, http.StatusOK, res)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func handleResourceAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}
	resourceID := parts[3]

	mu.RLock()
	res, exists := resources[resourceID]
	mu.RUnlock()

	if !exists {
		http.Error(w, `{"error": "Resource not found"}`, http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet && len(parts) == 4 {
		respondJSON(w, http.StatusOK, res)
		return
	}

	if r.Method == http.MethodPost && len(parts) == 5 && parts[4] == "power" {
		simulateChaos(w)

		var req PowerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Action != "on" && req.Action != "off" {
			http.Error(w, `{"error": "Invalid action"}`, http.StatusBadRequest)
			return
		}

		mu.Lock()
		if req.Action == "on" {
			res.Status = "running"
		} else {
			res.Status = "stopped"
		}
		mu.Unlock()

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"state":  req.Action,
		})
		return
	}

	http.NotFound(w, r)
}

func simulateChaos(w http.ResponseWriter) {
	if isChaos(slowRate) {
		time.Sleep(6 * time.Second)
	} else if isChaos(failureRate) {
		http.Error(w, `{"error": "Upstream service unavailable"}`, http.StatusInternalServerError)
		panic(http.ErrAbortHandler)
	}
}

func isChaos(percentage int) bool {
	return rand.Intn(100) < percentage
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
