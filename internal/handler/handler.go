package handler

import (
	"encoding/json"
	"net/http"
	"payment_midtrans/dto"
	"payment_midtrans/internal/usecase"
	"payment_midtrans/utils"

	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentHandler struct {
	payUC usecase.PaymentUsecase
}

func NewPaymentHandler(payUC usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{payUC}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {

	//now let req is empty beceuse req already filled in usecase for example
	var req snap.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.payUC.CreatePayment(&req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJson(w, http.StatusOK, response)
}

func (h *PaymentHandler) PaymentNotification(w http.ResponseWriter, r *http.Request) {

	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.payUC.PaymentNotificationStatus(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJson(w, http.StatusOK, response)
}

func (h *PaymentHandler) RefundTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.RefundCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	response, err := h.payUC.Refund(&req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJson(w, http.StatusOK, response)
}

func (h *PaymentHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req dto.PaymentSubcriptionCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response, err := h.payUC.CreateSubscription(&req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}

	utils.WriteJson(w, http.StatusOK, response)
}

func (h *PaymentHandler) SubscriptionNotification(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
	}
	respose, err := h.payUC.SubscriptionNotification(req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJson(w, http.StatusOK, respose)
}

func (h *PaymentHandler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	paramsId := params["subId"]

	response, err := h.payUC.CancelSubscription(paramsId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{
		"message": response,
	})
}
