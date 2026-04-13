package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"pbf-bridge/internal/domain"
)

type PrintHandler struct {
	shippingUseCase domain.ShippingUseCase
	identityUseCase domain.IdentityUseCase
}

func NewPrintHandler(suc domain.ShippingUseCase, iuc domain.IdentityUseCase) *PrintHandler {
	return &PrintHandler{
		shippingUseCase: suc,
		identityUseCase: iuc,
	}
}

func setupCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// 1. HANDLER FOR SHIPPING LABEL
func (h *PrintHandler) PrintShipping(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var payload domain.PrintShippingPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("Failed to decode shipping payload", "error", err.Error())
		http.Error(w, `{"error": "invalid payload"}`, http.StatusBadRequest)
		return
	}

	slog.Info("Received Shipping Request", "customer", payload.Recipient.Customer, "total_box", payload.TotalBox)

	if err := h.shippingUseCase.ProcessShippingLabels(payload, false); err != nil {
		slog.Error("Failed to process shipping label", "error", err.Error())
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Shipping labels printed"}`))
}

// 2. HANDLER FOR IDENTITY LABEL (PRODUCT LABEL)
func (h *PrintHandler) PrintIdentity(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var payload domain.PrintIdentityPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("Failed to decode identity payload", "error", err.Error())
		http.Error(w, `{"error": "invalid payload"}`, http.StatusBadRequest)
		return
	}

	slog.Info("Menerima request cetak Identity", "total_labels", len(payload.Identities))

	if err := h.identityUseCase.ProcessIdentityLabels(payload, false); err != nil {
		slog.Error("Failed to process identity label", "error", err.Error())
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Identity labels printed"}`))
}
