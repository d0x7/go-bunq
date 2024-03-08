package bunq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/d0x7/go-bunq/model"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type sessionServerService service

func (s *sessionServerService) create() (*model.ResponseSessionServer, error) {
	bodyStruct := model.RequestSessionServer{
		Secret: s.client.apiKey,
	}
	bodyRaw, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not marshal body")
	}

	r, err := http.NewRequest(
		http.MethodPost,
		s.client.formatRequestURL(endpointSessionServerCreate),
		bytes.NewBuffer(bodyRaw),
	)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: could not create request for session-server")
	}

	res, err := s.client.do(r)
	if err != nil {
		return nil, errors.Wrap(err, "bunq: request to session-server failed")
	}

	var resSessionServer model.ResponseSessionServer
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resSessionServer)
	if err != nil {
		return &resSessionServer, errors.Wrap(err, "bunq: decoding response into strcut failed")
	}

	resSessionServerProp := createProperSessionServerResponse(&resSessionServer)
	s.updateClient(resSessionServerProp)

	return resSessionServerProp, err
}

func createProperSessionServerResponse(r *model.ResponseSessionServer) *model.ResponseSessionServer {
	return &model.ResponseSessionServer{
		Response: []model.SessionServer{
			{
				ID:          r.Response[0].ID,
				Token:       r.Response[1].Token,
				UserPerson:  r.Response[2].UserPerson,
				UserCompany: r.Response[2].UserCompany,
				UserAPIKey:  r.Response[2].UserAPIKey,
			},
		},
	}
}

func (s *sessionServerService) updateClient(r *model.ResponseSessionServer) {
	s.updateClientToken(r)
	s.client.updateUserFlag()
}

func (s *sessionServerService) updateClientToken(r *model.ResponseSessionServer) {
	s.client.sessionServerContext = &r.Response[0]

	s.client.tokenMutex.Lock()
	defer s.client.tokenMutex.Unlock()
	s.client.token = &s.client.sessionServerContext.Token.Token

	if s.client.Debug {
		log.Printf("bunq: updating client token to session token %q", *s.client.token)
	}
}

func (s *sessionServerService) delete() error {
	url := s.client.formatRequestURL(fmt.Sprintf("session/%d", s.client.sessionServerContext.ID.ID))
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("bunq: could not create request for  %s", url))
	}

	res, err := s.client.Do(r)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("bunq: request to %s failed", url))
	}

	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("bunq: request to delete session resulted in response code: %d", res.StatusCode))
	}

	return err
}
