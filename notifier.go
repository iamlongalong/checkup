package checkup

import (
	"encoding/json"
	"fmt"

	"github.com/iamlongalong/checkup/notifier/discord"
	"github.com/iamlongalong/checkup/notifier/feishu"
	"github.com/iamlongalong/checkup/notifier/mail"
	"github.com/iamlongalong/checkup/notifier/mailgun"
	"github.com/iamlongalong/checkup/notifier/pushover"
	"github.com/iamlongalong/checkup/notifier/slack"
)

func notifierDecode(typeName string, config json.RawMessage) (Notifier, error) {
	switch typeName {
	case mail.Type:
		return mail.New(config)
	case slack.Type:
		return slack.New(config)
	case feishu.Type:
		return feishu.New(config)
	case mailgun.Type:
		return mailgun.New(config)
	case pushover.Type:
		return pushover.New(config)
	case discord.Type:
		return discord.New(config)
	default:
		return nil, fmt.Errorf(errUnknownNotifierType, typeName)
	}
}
