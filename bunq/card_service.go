package bunq

import (
	"fmt"
	"github.com/d0x7/go-bunq/model"
	"net/http"
)

type cardService service

func (c *cardService) GetMasterCardAction(id, monetaryAccountID int) (*model.ResponseMasterCardActionGet, error) {
	userID, err := c.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := c.client.preformRequest(http.MethodGet, c.client.formatRequestURL(fmt.Sprintf(endpointMasterCardActionGet, userID, monetaryAccountID, id)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponseMasterCardActionGet

	return &resStruct, c.client.parseResponse(res, &resStruct)
}
