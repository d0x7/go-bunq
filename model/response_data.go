package model

type ResponseInstallation struct {
	Response []Installation
}

type ResponseError struct {
	Error []bunqError `json:"Error"`
}

type ResponseDeviceServer struct {
	Response []wrappedBunqID
}

type ResponseSessionServer struct {
	Response []SessionServer
}

type ResponseUserPerson struct {
	Response []struct {
		UserPerson userPerson
	}
}

type ResponseBunqID struct {
	Response []wrappedBunqID
}

// ResponseMonetaryAccountBankGet The monetary account bank response object.
type ResponseMonetaryAccountBankGet struct {
	Response []struct {
		MonetaryAccountBank MonetaryAccountBank `json:"MonetaryAccountBank"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

// ResponseMonetaryAccountSavingGet The monetary account savings response object.
type ResponseMonetaryAccountSavingGet struct {
	Response []struct {
		MonetaryAccountSaving MonetaryAccountSaving `json:"MonetaryAccountSavings"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

type ResponseDraftPaymentGet struct {
	Response []struct {
		DraftPayment draftPayment `json:"DraftPayment"`
	} `json:"Response"`
}

// ResponsePaymentGet The payment response data.
type ResponsePaymentGet struct {
	Response []struct {
		Payment Payment `json:"Payment"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

type ResponseMasterCardActionGet struct {
	Response []struct {
		MasterCardAction masterCardAction `json:"MasterCardAction"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

// ResponseScheduledPaymentsGet The scheduled payments response object.
type ResponseScheduledPaymentsGet struct {
	Response []struct {
		ScheduledPayment ScheduledPayment `json:"ScheduledPayment"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}

type bunqError struct {
	ErrorDescription           string `json:"error_description"`
	ErrorDescriptionTranslated string `json:"error_description_translated"`
}

type ResponseRequestResponsesGet struct {
	Response []struct {
		RequestResponse RequestResponse `json:"RequestResponse"`
	} `json:"Response"`
	Pagination Pagination `json:"Pagination"`
}
