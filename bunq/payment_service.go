package bunq

import (
	"encoding/json"
	"fmt"
	"github.com/d0x7/go-bunq/model"
	"net/http"

	"github.com/pkg/errors"
)

type paymentService service

func (p *paymentService) CreateDraftPayment(monetaryAccountID int, rBody model.RequestCreateDraftPayment) (*model.ResponseBunqID, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	bodyRaw, err := json.Marshal(rBody)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	return p.client.doCURequest(p.client.formatRequestURL(fmt.Sprintf(endpointDraftPaymentCreate, userID, monetaryAccountID)), bodyRaw, http.MethodPost)
}

func (p *paymentService) UpdateDraftPayment(id, monetaryAccountID int, rBody model.RequestUpdateDraftPayment) (*model.ResponseBunqID, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	bodyRaw, err := json.Marshal(rBody)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	return p.client.doCURequest(p.client.formatRequestURL(fmt.Sprintf(endpointDraftPaymentWithID, userID, monetaryAccountID, id)), bodyRaw, http.MethodPut)
}

func (p *paymentService) GetDraftPayment(id, monetaryAccountID int) (*model.ResponseDraftPaymentGet, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(fmt.Sprintf(endpointDraftPaymentWithID, userID, monetaryAccountID, id)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponseDraftPaymentGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

// GetPayment returns a specific payment for a given account
func (p *paymentService) GetPayment(monetaryAccountID int, paymentID int) (*model.ResponsePaymentGet, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, errors.Wrap(err, "bunq: payment service: could not determine user id")
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(fmt.Sprintf(endpointPaymentGetWithID, userID, monetaryAccountID, paymentID)), nil)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponsePaymentGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

// GetAllPayment returns all the payments for a given account
func (p *paymentService) GetAllPayment(monetaryAccountID int, params ...model.QueryParam) (*model.ResponsePaymentGet, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, errors.Wrap(err, "bunq: payment service: could not determine user id")
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(fmt.Sprintf(endpointPaymentGet, userID, monetaryAccountID)), nil, params...)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponsePaymentGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

// GetAllOlderPayment calls the older url from the Pagination
func (p *paymentService) GetAllOlderPayment(pagi model.Pagination) (*model.ResponsePaymentGet, error) {
	if pagi.OlderURL == "" {
		return nil, nil
	}

	res, err := p.client.preformRequest(http.MethodGet, p.client.formatRequestURL(pagi.OlderURL[len("/v1/"):]), nil)
	if err != nil {
		return nil, err
	}

	var resStruct model.ResponsePaymentGet

	return &resStruct, p.client.parseResponse(res, &resStruct)
}

func (p *paymentService) CreatePaymentBatch(monetaryAccountID int, create model.PaymentBatchCreate) (*model.ResponseBunqID, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	bodyRaw, err := json.Marshal(create)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	return p.client.doCURequest(p.client.formatRequestURL(fmt.Sprintf(endpointPaymentBatchCreate, userID, monetaryAccountID)), bodyRaw, http.MethodPost)
}

func (p *paymentService) CreatePayment(monetaryAccountID int, create model.PaymentCreate) (*model.ResponseBunqID, error) {
	userID, err := p.client.GetUserID()
	if err != nil {
		return nil, err
	}

	bodyRaw, err := json.Marshal(create)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	return p.client.doCURequest(p.client.formatRequestURL(fmt.Sprintf(endpointPaymentCreate, userID, monetaryAccountID)), bodyRaw, http.MethodPost)
}
