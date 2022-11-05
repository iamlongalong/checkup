package checkup

import (
	"encoding/json"
	"fmt"

	"github.com/iamlongalong/checkup/check/dns"
	"github.com/iamlongalong/checkup/check/exec"
	"github.com/iamlongalong/checkup/check/http"
	"github.com/iamlongalong/checkup/check/tcp"
	"github.com/iamlongalong/checkup/check/tls"
)

func checkerDecode(typeName string, config json.RawMessage) (Checker, error) {
	switch typeName {
	case dns.Type:
		return dns.New(config)
	case exec.Type:
		return exec.New(config)
	case http.Type:
		return http.New(config)
	case tcp.Type:
		return tcp.New(config)
	case tls.Type:
		return tls.New(config)
	default:
		return nil, fmt.Errorf(errUnknownCheckerType, typeName)
	}
}
