package main

import (
	"log"
	"net/http"
	"os"
	"payment_midtrans/internal/handler"
	"payment_midtrans/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {

	// Setup routes
	payH := handler.NewPaymentHandler(usecase.NewPaymentUsecase())

	r := mux.NewRouter()
	r.HandleFunc("/payment", payH.CreatePayment).Methods(http.MethodPost)
	r.HandleFunc("/payment/notification", payH.PaymentNotification).Methods(http.MethodPost)
	r.HandleFunc("/payment/refund", payH.RefundTransaction).Methods(http.MethodPost)

	r.HandleFunc("/payment/subcription", payH.CreateSubscription).Methods(http.MethodPost)
	r.HandleFunc("/payment/subcription/notification", payH.SubscriptionNotification).Methods(http.MethodPost)
	r.HandleFunc("/payment/subcription/cancel", payH.CancelSubscription).Methods(http.MethodPost)
	//the port must be same on ngrok port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
