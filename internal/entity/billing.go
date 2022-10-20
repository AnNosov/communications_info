package entity

type BillingData struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"recurring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}
