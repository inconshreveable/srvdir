package server

import (
	"github.com/inconshreveable/go-tunnel/proto"
	"github.com/inconshreveable/go-tunnel/server"
	"os"
)

type SessionHooks struct {
	reg     *RegistrationHooks
	metrics *MetricsHooks
}

func NewSessionHooks() *SessionHooks {
	hooks := new(SessionHooks)

	// enable metrics if creds provided
	keenApiKey := os.Getenv("KEEN_API_KEY")
	keenProjectToken := os.Getenv("KEEN_PROJECT_TOKEN")
	if keenApiKey != "" && keenProjectToken != "" {
		hooks.metrics = NewMetricsHooks(keenApiKey, keenProjectToken)
	}

	// enable registration if URLs provided
	regAuthURL := os.Getenv("REG_AUTH_URL")
	regBindURL := os.Getenv("REG_BIND_URL")
	if regAuthURL != "" && regBindURL != "" {
		hooks.reg = NewRegistrationHooks(regAuthURL, regBindURL)
        hooks.reg.Info("Registration hooks enabled for URLs: Auth: %v, Bind: %v", regAuthURL, regBindURL)
	}

	return hooks
}

func (h *SessionHooks) OnAuth(sess *server.Session, auth *proto.Auth) error {
	if h.reg != nil {
		if err := h.reg.OnAuth(sess, auth); err != nil {
			return err
		}
	}

	return nil
}

func (h *SessionHooks) OnBind(sess *server.Session, bind *proto.Bind) error {
	if h.reg != nil {
		if err := h.reg.OnBind(sess, bind); err != nil {
			return err
		}
	}

	return nil
}

func (h *SessionHooks) OnClose(sess *server.Session) error {
	if h.metrics != nil {
		if err := h.metrics.OnClose(sess); err != nil {
			return err
		}
	}

	return nil
}
