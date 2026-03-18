package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tinylynx/internal/service"
)

type LinkHandler struct{}

func NewLinkHandler() *LinkHandler {
	return &LinkHandler{}
}

func (h *LinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OriginalLink string `json:"original_link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OriginalLink == "" {
		http.Error(w, "original_link is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	link, err := service.GetByOriginalLink(ctx, req.OriginalLink)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create link: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"short_code": link.ShortCode,
		"original_link": link.OriginalLink,
		"created_at": link.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	link, err := service.FindByShortCode(ctx, shortCode)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	// Record analytics
	ip := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	referrer := r.Header.Get("Referer")
	country := r.Header.Get("X-Forwarded-For")
	device := detectDevice(r.UserAgent())
	browser := detectBrowser(r.UserAgent())
	platform := detectPlatform(r.UserAgent())

	if country == "" {
		country = "unknown"
	}

	if err := service.RecordAnalytics(ctx, link.ID, ip, userAgent, referrer, country, device, browser, platform); err != nil {
		// Log error but don't fail the redirect
		fmt.Printf("Failed to record analytics: %v\n", err)
	}

	http.Redirect(w, r, link.OriginalLink, http.StatusTemporaryRedirect)
}

func (h *LinkHandler) GetLinkStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Query().Get("code")
	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	link, err := service.FindByShortCode(ctx, shortCode)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	stats, err := service.GetLinkStats(ctx, link.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *LinkHandler) GetLinkAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Query().Get("code")
	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	link, err := service.FindByShortCode(ctx, shortCode)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	analytics, err := service.GetLinkAnalytics(ctx, link.ID, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get analytics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

func detectDevice(userAgent string) string {
	if strings.Contains(userAgent, "Mobile") || strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPhone") {
		return "mobile"
	}
	if strings.Contains(userAgent, "iPad") {
		return "tablet"
	}
	return "desktop"
}

func detectBrowser(userAgent string) string {
	if strings.Contains(userAgent, "Chrome") {
		return "chrome"
	}
	if strings.Contains(userAgent, "Firefox") {
		return "firefox"
	}
	if strings.Contains(userAgent, "Safari") {
		return "safari"
	}
	if strings.Contains(userAgent, "Edge") {
		return "edge"
	}
	return "unknown"
}

func detectPlatform(userAgent string) string {
	if strings.Contains(userAgent, "Windows") {
		return "windows"
	}
	if strings.Contains(userAgent, "Mac") {
		return "macos"
	}
	if strings.Contains(userAgent, "Linux") {
		return "linux"
	}
	if strings.Contains(userAgent, "Android") {
		return "android"
	}
	if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") {
		return "ios"
	}
	return "unknown"
}