package bunq

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/d0x7/go-bunq/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

const errorRequestToUnMockedHTTPMethod = "request made to an un mocked http method. url: %q method: %q"

func createBunqFakeHandler(t *testing.T) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path[4:] {
		case "installation":
			err := json.NewEncoder(w).Encode(getInstallationResponse(t))
			if err != nil {
				t.Fatal(err)
			}
		case "device-server":
			sendResponseWithSignature(t, w, http.StatusOK, getDeviceServerResponse(t))
		case "session-server":
			sendResponseWithSignature(t, w, http.StatusOK, getSessionServerResponse(t))
		case "user-person/6084":
			switch r.Method {
			case http.MethodGet:
				sendResponseWithSignature(t, w, http.StatusOK, getUserPersonGetResponse(t))
			case http.MethodPut:
				sendResponseWithSignature(t, w, http.StatusOK, getGenericIDResponse(t))
			default:
				t.Errorf(errorRequestToUnMockedHTTPMethod, r.URL, r.Method)
			}
		case "user/6084/monetary-account/9512/draft-payment":
			switch r.Method {
			case http.MethodPost:
				sendResponseWithSignature(t, w, http.StatusOK, getGenericIDResponse(t))
			default:
				t.Errorf(errorRequestToUnMockedHTTPMethod, r.URL, r.Method)
			}
		case "user/6084/monetary-account-bank", "user/6084/monetary-account-bank/9601":
			sendResponseWithSignature(t, w, http.StatusOK, getMonetaryAccountBankGet(t))
		case "user/6084/monetary-account-savings", "user/6084/monetary-account-savings/9601":
			sendResponseWithSignature(t, w, http.StatusOK, getMonetaryAccountSavings(t))
		case "user/6084/monetary-account/9618/draft-payment", "user/6084/monetary-account/9618/draft-payment/6292":
			switch r.Method {
			case http.MethodPost, http.MethodPut:
				sendResponseWithSignature(t, w, http.StatusOK, getGenericIDResponse(t))
			case http.MethodGet:
				sendResponseWithSignature(t, w, http.StatusOK, getDraftPaymentGet(t))
			default:
				t.Errorf(errorRequestToUnMockedHTTPMethod, r.URL, r.Method)
			}

		case "user/6084/monetary-account/9520/mastercard-action/324":
			sendResponseWithSignature(t, w, http.StatusOK, getMasterCardActionGet(t))
		case "user/6084/monetary-account/10111/payment", "user/7082/monetary-account/10111/payment", "user/6084/monetary-account/10111/payment/1":
			sendResponseWithSignature(t, w, http.StatusOK, getPaymentGet(t))
		case "user/6084/monetary-account/9601/schedule-payment":
			sendResponseWithSignature(t, w, http.StatusOK, getScheduledPaymentGet(t))
		case "user/6084/monetary-account/9999/request-response":
			sendResponseWithSignature(t, w, http.StatusOK, getRequestResponseGet(t))
		case "attachment-public/f9a1a89a-fdc1-4de5-89d5-e477cccd22c4/content":
			sendResponseWithSignature(t, w, http.StatusOK, getPaymentGet(t))
		case "/v1/session/133912", "v1/session/133912", "session/133912":
			sendResponseWithSignature(t, w, http.StatusOK, getSessionServerResponse(t))
		default:
			t.Errorf("requst made to an un mocked url: %v", r.URL)
		}
	})
}

func createBunqFakeHandlerWithError(t *testing.T, endpointToError string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == endpointToError {
			sendResponseWithSignature(t, w, http.StatusTeapot, getErrorResponse(t))
		} else {
			createBunqFakeHandler(t)(w, r)
		}
	})
}

func createClientWithFakeServer(t *testing.T) (*Client, *httptest.Server, context.CancelFunc) {
	fakeServer := httptest.NewServer(createBunqFakeHandler(t))

	key, err := CreateNewKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := NewClient(ctx, fmt.Sprintf("%s/v1/", fakeServer.URL), key, "", "", CurrentIP)

	return c, fakeServer, cancel
}

func getInstallationResponse(t *testing.T) *model.ResponseInstallation {
	var obj model.ResponseInstallation
	res := createResponseStruct(t, formatFilePathByName("installation_response"), &obj)

	return res.(*model.ResponseInstallation)
}

func getDeviceServerResponse(t *testing.T) *model.ResponseDeviceServer {
	var obj model.ResponseDeviceServer
	res := createResponseStruct(t, formatFilePathByName("device_server_response"), &obj)

	return res.(*model.ResponseDeviceServer)
}

func getSessionServerResponse(t *testing.T) *model.ResponseSessionServer {
	var obj model.ResponseSessionServer
	res := createResponseStruct(t, formatFilePathByName("session_server_response"), &obj)

	return res.(*model.ResponseSessionServer)
}

func getUserPersonGetResponse(t *testing.T) *model.ResponseUserPerson {
	var obj model.ResponseUserPerson
	res := createResponseStruct(t, formatFilePathByName("user_person_get_response"), &obj)

	return res.(*model.ResponseUserPerson)
}

func getMonetaryAccountBankGet(t *testing.T) *model.ResponseMonetaryAccountBankGet {
	var obj model.ResponseMonetaryAccountBankGet
	res := createResponseStruct(t, formatFilePathByName("monetary_account_bank_listing_response"), &obj)

	return res.(*model.ResponseMonetaryAccountBankGet)
}

func getMonetaryAccountSavings(t *testing.T) *model.ResponseMonetaryAccountSavingGet {
	var obj model.ResponseMonetaryAccountSavingGet
	res := createResponseStruct(t, formatFilePathByName("monetary_account_savings_response_get"), &obj)

	return res.(*model.ResponseMonetaryAccountSavingGet)
}

func getGenericIDResponse(t *testing.T) *model.ResponseBunqID {
	var obj model.ResponseBunqID
	res := createResponseStruct(t, formatFilePathByName("generic_id_response"), &obj)

	return res.(*model.ResponseBunqID)
}

func getDraftPaymentGet(t *testing.T) *model.ResponseDraftPaymentGet {
	var obj model.ResponseDraftPaymentGet
	res := createResponseStruct(t, formatFilePathByName("draft_payment_get_response"), &obj)

	return res.(*model.ResponseDraftPaymentGet)
}

func getMasterCardActionGet(t *testing.T) *model.ResponseMasterCardActionGet {
	var obj model.ResponseMasterCardActionGet
	res := createResponseStruct(t, formatFilePathByName("master_card_action_get_response"), &obj)

	return res.(*model.ResponseMasterCardActionGet)
}

func getPaymentGet(t *testing.T) *model.ResponsePaymentGet {
	var obj model.ResponsePaymentGet
	res := createResponseStruct(t, formatFilePathByName("payment_get_response"), &obj)

	return res.(*model.ResponsePaymentGet)
}

func getScheduledPaymentGet(t *testing.T) *model.ResponseScheduledPaymentsGet {
	var obj model.ResponseScheduledPaymentsGet
	res := createResponseStruct(t, formatFilePathByName("schedule_payment_response"), &obj)

	return res.(*model.ResponseScheduledPaymentsGet)
}

func getRequestResponseGet(t *testing.T) *model.ResponseRequestResponsesGet {
	var obj model.ResponseRequestResponsesGet
	res := createResponseStruct(t, formatFilePathByName("request_response_response"), &obj)

	return res.(*model.ResponseRequestResponsesGet)
}

func getErrorResponse(t *testing.T) *model.ResponseError {
	var obj model.ResponseError
	res := createResponseStruct(t, formatFilePathByName("error_response"), &obj)

	return res.(*model.ResponseError)
}

func createResponseStruct(t *testing.T, path string, obj interface{}) interface{} {
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(obj)
	if err != nil {
		t.Fatal(err)
	}

	return obj
}

func formatFilePathByName(fileName string) string {
	return fmt.Sprintf("../testdata/bunq/%s.json", fileName)
}

func sendResponseWithSignature(t *testing.T, w http.ResponseWriter, resCode int, body interface{}) {
	b, _ := json.Marshal(body)
	stringToSign := fmt.Sprintf("%s", b)

	h := sha256.New()
	_, _ = h.Write([]byte(stringToSign))

	signature, _ := rsa.SignPKCS1v15(rand.Reader, loadPrivateKey(), crypto.SHA256, h.Sum(nil))
	w.Header().Set("X-Bunq-Server-Signature", base64.StdEncoding.EncodeToString(signature))
	w.WriteHeader(resCode)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
}

var loadPrivateKeyOnce sync.Once
var privateKey *rsa.PrivateKey

func loadPrivateKey() *rsa.PrivateKey {
	loadPrivateKeyOnce.Do(func() {
		k, _ := ioutil.ReadFile("../testdata/bunq/private.key")

		block, _ := pem.Decode(k)
		privateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	})

	return privateKey
}
