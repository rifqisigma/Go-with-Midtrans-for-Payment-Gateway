package dto

type PaymentCreateResponse struct {
	Token       string `json:"snap_token"`
	URL         string `json:"redirect_url"`
	OrderId     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

type PaymentNotifResponse struct {
	OrderId     string `json:"order_id"`
	Status      string `json:"transaction_status"`
	FraudStatus string `json:"fraud_status"`
}

type RefundCreate struct {
	OrderId   string `json:"order_id"`
	Amount    int64  `json:"amount"`
	ReasonWHY string `json:"reason"`
}

type PaymentSubcriptionCreate struct {
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	Amount        int64  `json:"amount"`
	TokenID       string `json:"token_id"`
	Interval      string `json:"interval"`
	IntervalCount int    `json:"interval_count"`
}

type SubscriptionResponse struct {
	UserID          string `json:"user_id"`
	MidtransSubID   string `json:"midtrans_midtrans_sub_id"`
	Status          string `json:"status"`
	NextBillingDate string `json:"next_billing_date"`
}
