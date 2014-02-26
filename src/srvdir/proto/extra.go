package extra

import (
	"github.com/inconshreveable/go-tunnel/proto"
	"runtime"
)

const Version = "1"

type AuthExtra struct {
	AuthToken            string
	OS                   string
	ClientVersion        string
	ExtraProtocolVersion string
}

func NewAuthExtra(token string, clientVersion string) *AuthExtra {
	return &AuthExtra{
		AuthToken:            token,
		OS:                   runtime.GOOS,
		ExtraProtocolVersion: Version,
		ClientVersion:        clientVersion,
	}
}

func UnpackAuthExtra(extra interface{}) (authExtra *AuthExtra, err error) {
	authExtra = new(AuthExtra)
	err = proto.UnpackInterfaceField(extra, authExtra)
	return
}
