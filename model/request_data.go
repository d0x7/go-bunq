package model

type RequestInstallation struct {
	ClientPublicKey string `json:"client_public_key"`
}

type RequestDeviceServer struct {
	Description  string   `json:"description"`
	Secret       string   `json:"secret"`
	PermittedIps []string `json:"permitted_ips"`
}

type RequestSessionServer struct {
	Secret string `json:"secret"`
}

type RequestUserPersonPut struct {
	NotificationFilters []NotificationFilter `json:"notification_filters,omitempty"`
}

type RequestCreateDraftPayment struct {
	Entries                 []DraftPaymentEntryCreate `json:"entries"`
	NumberOfRequiredAccepts *int                      `json:"number_of_required_accepts,omitempty"`
}

type RequestUpdateDraftPayment struct {
	RequestCreateDraftPayment
	UpdatedTimestamp string `json:"previous_updated_timestamp"`
	status           *string
}

type DraftPaymentEntryCreate struct {
	Amount            Amount  `json:"amount,omitempty"`
	CounterpartyAlias Pointer `json:"counterparty_alias,omitempty"`
	Description       string  `json:"description,omitempty"`
	MerchantReference *string `json:"merchant_reference,omitempty"`
}

type PaymentBatchCreate struct {
	Payments []PaymentCreate `json:"payments"`
}

type PaymentCreate struct {
	Amount            Amount  `json:"amount"`
	CounterpartyAlias Pointer `json:"counterparty_alias"`
	Description       string  `json:"description"`
	AllowBunqto       bool    `json:"allow_bunqto"`
}
