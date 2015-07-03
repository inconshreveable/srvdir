package server

import (
	"github.com/inconshreveable/go-keen"
	"github.com/inconshreveable/go-tunnel/log"
	"github.com/inconshreveable/go-tunnel/proto"
	"github.com/inconshreveable/go-tunnel/server"
	//"github.com/inconshreveable/go-tunnel/conn"
	srvdir_proto "srvdir/proto"
	"time"
)

const flushInterval = 10 * time.Second

type MetricsHooks struct {
	log.Logger
	client *keen.BatchClient
}

func NewMetricsHooks(keenApiKey, keenProjectToken string) *MetricsHooks {
	client := &keen.Client{WriteKey: keenApiKey, ProjectID: keenProjectToken}

	return &MetricsHooks{
		Logger: log.NewTaggedLogger("hooks", "metrics"),
		client: keen.NewBatchClient(client, flushInterval),
	}
}

type sessionCloseEvent struct {
	Keen keen.KeenProperties `json: "keen"`

	OS       string
	ClientId string
	User     string
	Version  string
	Reason   string
	Duration float64
}

type tunnelCloseEvent struct {
	Keen keen.KeenProperties `json:"keen"`

	OS        string
	ClientId  string
	User      string
	Version   string
	Duration  float64
	Protocol  string
	Url       string
	Auth      bool
	Subdomain bool
	Hostname  bool
}

type connectionCloseEvent struct {
	Keen keen.KeenProperties `json:"keen"`

	OS             string
	ClientId       string
	User           string
	Version        string
	Protocol       string
	Url            string
	Auth           bool
	Subdomain      bool
	Hostname       bool
	Duration       float64
	TunnelDuration float64
	BytesIn        int64
	BytesOut       int64
}

func (h *MetricsHooks) OnClose(sess *server.Session) error {
	_, authExtra, err := extractSession(sess)
	if err != nil {
		return err
	}

	err = h.client.AddEvent("CloseSession", &sessionCloseEvent{
		Keen: keen.KeenProperties{
			Timestamp: keen.Timestamp(sess.Start()),
		},
		OS:       authExtra.OS,
		ClientId: sess.Id(),
		User:     authExtra.AuthToken,
		Version:  authExtra.ClientVersion,
		Reason:   "", // XXX
		Duration: time.Since(sess.Start()).Seconds(),
	})

	if err != nil {
		h.Error("Failed to record CloseSession event: %v", err)
	}

	return nil
}

/*
func (h *MetricsHooks) OnTunnelClose(tunnel *server.Tunnel) error {
    sess := tunnel.Session()
    if auth, authExtra, err := extractSession(sess); err != nil {
        return
    }

    if bind, opts, err := extractTunnel(tunnel); err != nil {
        return
    }

    err := client.AddEvent("CloseTunnel", &tunnelCloseEvent{
        Keen: keen.KeenProperties {
            Timestamp: keen.Timestamp(tunnel.Start()),
        },

        OS: authExtra.OS,
        ClientId: sess.Id(),
        User: authExtra.User,
        Version: authExtra.Version,
        Duration: time.Since(tunnel.Start()),
        Protocol: bind.Protocol,
        Url: tunnel.Url(),
        Auth: opts.Auth != "",
        Subdomain: opts.Subdomain != "",
        Hostname: opts.Hostname != "",
    })
}

func (h *MetricsHooks) OnConnectionClose(tunnel *server.Tunnel, c conn.Conn, duration time.Duration, bytesIn, bytesOut int64) error {
    sess := tunnel.Session()
    if auth, authExtra, err := extractSession(sess); err != nil {
        return
    }

    if bind, opts, err := extractTunnel(tunnel); err != nil {
        return
    }

    err := client.AddEvent("CloseConnection", &connectionCloseEvent{
        Keen: keen.KeenProperties {
            Timestamp: keen.Timestamp(tunnel.Start()),
        },
        OS: authExtra.OS,
        ClientId: sess.Id(),
        User: authExtra.User,
        Version: authExtra.Version,
        Protocol: bind.Protocol,
        Url: tunnel.Url(),
        Auth: opts.Auth != "",
        Subdomain: opts.Subdomain != "",
        Hostname: opts.Hostname != "",
        Duration: duration,
        TunnelDuration: time.Since(tunnel.Start()),
        BytesIn: bytesIn,
        BytesOut: bytesOut,
    })
}
*/

func extractSession(s *server.Session) (*proto.Auth, *srvdir_proto.AuthExtra, error) {
	var auth proto.Auth
	if err := proto.UnpackInterfaceField(s.Auth(), &auth); err != nil {
		return nil, nil, err
	}

	var authExtra srvdir_proto.AuthExtra
	if err := proto.UnpackInterfaceField(s.Auth().Extra, &authExtra); err != nil {
		return nil, nil, err
	}

	return &auth, &authExtra, nil
}

/*
func extraTunnel(t *server.Tunnel) (*proto.Bind, *proto.HTTPOptions, error) {
    bind := t.Bind()
    var opts proto.HTTPOptions
    if err := proto.UnpackInterfaceField(bind.Options, &opts); err != nil {
        return nil, nil, err
    }

    return bind, &opts, nil
}
*/
