package usecase

import (
	"fmt"
	"os"
	"payment_midtrans/dto"
	"payment_midtrans/utils"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentUsecase interface {
	CreatePayment(req *snap.Request) (*dto.PaymentCreateResponse, error)
	PaymentNotificationStatus(payload map[string]interface{}) (*dto.PaymentNotifResponse, error)
	Refund(req *dto.RefundCreate) (*coreapi.RefundResponse, error)
	CreateSubscription(req *dto.PaymentSubcriptionCreate) (*dto.SubscriptionResponse, error)
	SubscriptionNotification(payload map[string]interface{}) (*dto.SubscriptionResponse, error)
	CancelSubscription(subId string) (string, error)
	CheckSubscription(subscriptionID string) (*coreapi.StatusSubscriptionResponse, error)
}

type paymentUsecase struct {
}

func NewPaymentUsecase() PaymentUsecase {
	return &paymentUsecase{}
}

func (u *paymentUsecase) CreatePayment(req *snap.Request) (*dto.PaymentCreateResponse, error) {
	s := snap.Client{}

	fmt.Println(os.Getenv("MIDTRANS_SERVER_KEY"))
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	//for example
	req = &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "ORDER-" + time.Now().Format("20060102150405"),
			GrossAmt: 100000,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "John",
			LName: "Doe",
			Email: "john.doe@example.com",
			Phone: "08123456789",
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "item01",
				Name:  "T-Shirt",
				Price: 50000,
				Qty:   2,
			},
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeShopeepay,
			snap.PaymentTypeGopay,
		},
		Expiry: &snap.ExpiryDetails{
			StartTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
			Unit:      "minute",
			Duration:  15,
		},
	}

	snapResp, err := s.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentCreateResponse{
		OrderId:     req.TransactionDetails.OrderID,
		URL:         snapResp.RedirectURL,
		Token:       snapResp.Token,
		GrossAmount: req.TransactionDetails.GrossAmt,
	}, nil
}

func (u *paymentUsecase) PaymentNotificationStatus(payload map[string]interface{}) (*dto.PaymentNotifResponse, error) {
	var core coreapi.Client
	core.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	//for postman testing purpose
	orderID, _ := payload["order_id"].(string)

	// statusCode, _ := payload["status_code"].(string)
	// grossAmount, _ := payload["gross_amount"].(string)
	// receivedSignature, _ := payload["signature_key"].(string)

	// if orderID == "" || statusCode == "" || grossAmount == "" || receivedSignature == "" {
	// 	return nil, errors.New("missing required fields")
	// }

	// expectedSig := utils.GenerateSignature(orderID, statusCode, grossAmount, os.Getenv("MIDTRANS_SERVER_KEY"))
	// if receivedSignature != expectedSig {
	// 	return nil, errors.New("failed signature")
	// }

	transactionStatusResp, errResp := core.CheckTransaction(orderID)
	if errResp != nil {
		return nil, errResp
	}

	switch transactionStatusResp.TransactionStatus {
	case "capture":
		// credit card case
		if transactionStatusResp.FraudStatus == "accept" {
			// payment done UPDATE DB

		} else if transactionStatusResp.FraudStatus == "challenge" {
			if err := utils.VoidTransaction(orderID); err != nil {

				//update void_failed
				return nil, err
			}
			// update db status voided(dibatalkan)
		}
	case "settlement":
		// payment done UPDATE DB
	case "pending":

	case "deny", "cancel", "expire":

	default:
	}

	return &dto.PaymentNotifResponse{
		OrderId:     orderID,
		Status:      transactionStatusResp.TransactionStatus,
		FraudStatus: transactionStatusResp.FraudStatus,
	}, nil
}

func (u *paymentUsecase) Refund(req *dto.RefundCreate) (*coreapi.RefundResponse, error) {

	var core coreapi.Client
	core.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	refundKey := fmt.Sprint("refund-%S", time.Now().Unix())
	refundReq := &coreapi.RefundReq{
		RefundKey: refundKey,
		Amount:    req.Amount,
		Reason:    req.ReasonWHY,
	}

	refundResp, err := core.RefundTransaction(req.OrderId, refundReq)
	if err != nil {
		return nil, err
	}

	return refundResp, nil

}

func (u *paymentUsecase) CreateSubscription(req *dto.PaymentSubcriptionCreate) (*dto.SubscriptionResponse, error) {

	core := coreapi.Client{}
	core.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	//sub
	subReq := &coreapi.SubscriptionReq{
		Name:        "john",
		Amount:      10000,
		Currency:    "IDR",
		PaymentType: "credit_card",
		Token:       "fake-token-123",
		Schedule: coreapi.ScheduleDetails{
			Interval:     1,
			IntervalUnit: "month",
			MaxInterval:  2,
			StartTime:    time.Now().Add(2 * time.Minute).Format("2006-01-02 15:04:05 -0700"),
		},
		Metadata: map[string]string{
			"deskripsi": "beli bulanan",
		},
	}

	subResp, err := core.CreateSubscription(subReq)
	if err != nil {
		return nil, err
	}

	sub := &dto.SubscriptionResponse{
		UserID:          "1",
		MidtransSubID:   subResp.ID,
		Status:          "active",
		NextBillingDate: subResp.Schedule.IntervalUnit,
	}

	return sub, nil
}

func (u *paymentUsecase) SubscriptionNotification(payload map[string]interface{}) (*dto.SubscriptionResponse, error) {

	subID := payload["subscription_id"].(string)
	status := payload["status"].(string)
	nextBilling, _ := payload["next_payment_date"].(string)

	//query for update status and get db to get response

	sub := &dto.SubscriptionResponse{
		// UserID:
		MidtransSubID:   subID,
		Status:          status,
		NextBillingDate: nextBilling,
	}

	return sub, nil

}

func (u *paymentUsecase) CancelSubscription(subId string) (string, error) {

	core := coreapi.Client{}
	core.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	resp, err := core.DisableSubscription(subId)
	if err != nil {
		return "", err
	}

	//update db status canceled

	return resp.StatusMessage, nil
}

func (u *paymentUsecase) CheckSubscription(subscriptionID string) (*coreapi.StatusSubscriptionResponse, error) {
	// Inisialisasi client
	var core coreapi.Client
	core.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	// Panggil GetSubscription
	subscriptionDetail, err := core.GetSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	return subscriptionDetail, nil
}
