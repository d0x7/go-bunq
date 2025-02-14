package bunq

import (
	"fmt"
	"github.com/d0x7/go-bunq/model"
	"net/http"

	"github.com/pkg/errors"
)

type requestResponseService service

// GetAllRequestResponses returns all request responses for a given account
func (p *requestResponseService) GetAllRequestResponses(monetaryAccountID int, params ...model.QueryParam) (*model.ResponseRequestResponsesGet, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request-response service: could not determine user id")
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(fmt.Sprintf(endpointRequestResponsesGet, userID, monetaryAccountID)), nil, params...)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponseRequestResponsesGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

// GetRequestResponse returns a specific request response for a given account
func (p *requestResponseService) GetRequestResponse(monetaryAccountID int, requestResponseID int) (*model.ResponseRequestResponsesGet, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(fmt.Sprintf(endpointRequestResponsesGetWithID, userID, monetaryAccountID, requestResponseID)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponseRequestResponsesGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

// GetAllOlderRequestResponses calls the older url from the Pagination
func (p *requestResponseService) GetAllOlderRequestResponses(pagi model.Pagination) (*model.ResponseRequestResponsesGet, error) {
	if pagi.OlderURL == "" {
		return nil, nil
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(pagi.OlderURL[len("/v1/"):]), nil)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponseRequestResponsesGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}
