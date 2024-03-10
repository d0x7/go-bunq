package bunq

import (
	"fmt"
	"github.com/d0x7/go-bunq/model"
	"net/http"
)

type cardService service

func (c *cardService) GetAllMasterCardAction(monetaryAccountID int, params ...model.QueryParam) (*model.ResponseMasterCardActionGet, error) {
	userID, err := c.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := c.client.preformRequest(http.MethodGet, c.client.formatRequestURL(fmt.Sprintf(endpointMasterCardActionGet, userID, monetaryAccountID)), nil, params...)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponseMasterCardActionGet

	return &resStruct, c.client.parseResponse(res, &resStruct)
}

func (c *cardService) GetMasterCardAction(monetaryAccountID int, id int) (*model.ResponseMasterCardActionGet, error) {
	userID, err := c.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := c.client.preformRequest(http.MethodGet, c.client.formatRequestURL(fmt.Sprintf(endpointMasterCardActionGetWithID, userID, monetaryAccountID, id)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponseMasterCardActionGet

	return &resStruct, c.client.parseResponse(res, &resStruct)
}
