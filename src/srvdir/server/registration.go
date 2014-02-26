package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/inconshreveable/go-tunnel/log"
	"github.com/inconshreveable/go-tunnel/proto"
	"github.com/inconshreveable/go-tunnel/server"
	"io/ioutil"
	"net/http"
	srvdir_proto "srvdir/proto"
)

type RegistrationHooks struct {
	log.Logger

	onAuthURL string
	onBindURL string
}

func NewRegistrationHooks(onAuthURL, onBindURL string) *RegistrationHooks {
	return &RegistrationHooks{
		Logger:    log.NewTaggedLogger("hooks", "registration"),
		onAuthURL: onAuthURL,
		onBindURL: onBindURL,
	}
}

type AuthRegRequest struct {
	AuthExtra *srvdir_proto.AuthExtra
}

func (h *RegistrationHooks) OnAuth(sess *server.Session, auth *proto.Auth) error {
	authExtra, err := srvdir_proto.UnpackAuthExtra(auth.Extra)
	if err != nil {
		h.Warn("Authentication data is malformed: %v", err)
		return fmt.Errorf("Authentication data is malformed: %v", err)
	}

	payload, err := json.Marshal(&AuthRegRequest{authExtra})
	if err != nil {
		h.Error("Failed to build bind registration payload: %v", err)
		return fmt.Errorf("Internal error building registration payload")
	}

	resp, err := http.Post(h.onAuthURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		h.Error("Error consulting registration service: %v", err)
		return nil
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.Error("Error reading registration service response: %v", err)
		return nil
	}

	switch resp.StatusCode {
	case 400:
		var validationErr struct {
			Message string `json:message`
		}
		err = json.Unmarshal(body, &validationErr)
		if err != nil {
			h.Error("Failed to unmarshal registration service error response: %v", err)
			return nil
		}
		h.Info("Registration service failed request: %v", err)
		return fmt.Errorf("%v", validationErr.Message)

	case 200:
		// everything is OK

	default:
		h.Error("Registration service returned unhandled status code: %v", resp.StatusCode)
		return nil
	}

	return nil
}

type BindRegRequest struct {
	Auth *srvdir_proto.AuthExtra
	Bind *proto.Bind
}

func (h *RegistrationHooks) OnBind(sess *server.Session, bind *proto.Bind) error {
	authExtra, err := srvdir_proto.UnpackAuthExtra(sess.Auth())
	if err != nil {
		h.Warn("Authentication data is malformed: %v", err)
		return fmt.Errorf("Authentication data is malformed: %v", err)
	}

	payload, err := json.Marshal(&BindRegRequest{Auth: authExtra, Bind: bind})
	if err != nil {
		h.Error("Failed to build bind registration payload: %v", err)
		return fmt.Errorf("Internal error building registration payload")
	}

	resp, err := http.Post(h.onBindURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		h.Error("Error consulting registration service: %v", err)
		return nil
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.Error("Error reading registration service response: %v", err)
		return nil
	}

	switch resp.StatusCode {
	case 400:
		var validationErr struct {
			Message string `json:message`
		}
		err = json.Unmarshal(body, &validationErr)
		if err != nil {
			h.Error("Failed to unmarshal registration service error response: %v", err)
			return nil
		}
		h.Info("Registration service failed request: %v", err)
		return fmt.Errorf("%v", validationErr.Message)

	case 200:
		err = json.Unmarshal(body, bind)
		if err != nil {
			h.Error("Failed to unmarshal registration service response: %v", err)
            return nil
		}

	default:
		h.Error("Registration service returned unhandled status code: %v", resp.StatusCode)
		return nil
	}

	return nil
}
