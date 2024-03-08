package bunq

import (
	"fmt"
	"github.com/d0x7/go-bunq/model"
	"net/http"

	"github.com/pkg/errors"
)

type scheduledPaymentService service

func (sp *scheduledPaymentService) GetAllScheduledPayments(monetaryAccountID int) (*model.ResponseScheduledPaymentsGet, error) {
	userID, err := sp.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := sp.client.preformRequest(http.MethodGet, sp.client.formatRequestURL(fmt.Sprintf(endpointScheduledPaymentGet, userID, monetaryAccountID)), nil)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to get all scheduled payments failed")
	}

	var resSpGet model.ResponseScheduledPaymentsGet

	return &resSpGet, sp.client.parseResponse(res, &resSpGet)
}
