package feishu

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-lark/lark"
	"github.com/sourcegraph/checkup/types"
)

// Type should match the package name
const Type = "feishu"

// Notifier consist of all the sub components required to send E-mail notifications
type Notifier struct {
	mu *sync.Mutex

	Cfg struct {
		// From contains the e-mail address notifications are sent from
		WebHook string `json:"webhook"`

		// Alert Base Times
		BaseTimes int `json:"basetimes"` // 连续出错次数

		// Slient Time
		SlientTime string `json:"slienttime"` // 静默时间

		// Title
		Title string `json:"title"` // 消息的 title
	}

	failRecords map[string]*notifyRecord

	bot *lark.Bot
}

func newRecord(slientTime time.Duration, baseTimes int) *notifyRecord {
	return &notifyRecord{
		slientTime:    slientTime,
		baseFailTimes: baseTimes,
	}
}

func (r *notifyRecord) ShouldNotify() bool {
	if r.lastNotifyTime.Add(r.slientTime).After(time.Now()) { // 静默期内，不可执行
		r.failTimes++
		return false
	}

	if r.failTimes+1 >= r.baseFailTimes { // 已经超了，可以执行
		r.clear()
		return true
	}

	r.failTimes++

	return false
}

func (r *notifyRecord) clear() {
	r.lastNotifyTime = time.Now() // 恢复时间
	r.failTimes = 0               // 恢复出错次数
}

type notifyRecord struct {
	lastNotifyTime time.Time     // 上次通知事件
	slientTime     time.Duration // 静默时间

	failTimes     int // 当前出错次数
	baseFailTimes int // 基础触发 出错次数
}

// New creates a new Notifier instance based on json config
func New(config json.RawMessage) (Notifier, error) {
	var notifier Notifier
	err := json.Unmarshal(config, &notifier.Cfg)
	if err != nil {
		return notifier, err
	}

	bot := lark.NewNotificationBot(notifier.Cfg.WebHook)
	notifier.bot = bot
	notifier.mu = &sync.Mutex{}
	notifier.failRecords = map[string]*notifyRecord{}

	return notifier, err
}

// Type returns the notifier package name
func (Notifier) Type() string {
	return Type
}

// Notify implements notifier interface
func (m Notifier) Notify(results []types.Result) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Println("xxx")

	issues := []types.Result{}
	for _, result := range results {
		r, ok := m.failRecords[result.Title]
		if !ok {
			t, _ := time.ParseDuration(m.Cfg.SlientTime)

			r = newRecord(t, m.Cfg.BaseTimes)
			m.failRecords[result.Title] = r
		}

		if result.Healthy { // 健康，消除告警
			r.clear()
			continue
		}

		if r.ShouldNotify() {
			issues = append(issues, result)
		}
	}

	if len(issues) == 0 {
		return nil
	}

	msgstr := renderMsg(m.Cfg.Title+"\n", issues)

	msg := lark.OutcomingMessage{
		MsgType: lark.MsgText,
		Content: lark.MessageContent{
			Text: &lark.TextContent{
				Text: msgstr,
			},
		},
	}

	res, err := m.bot.PostNotificationV2(msg)
	if err != nil {
		return err
	}

	fmt.Println("res : ", res)

	return nil
}

func renderMsg(base string, results []types.Result) string {
	for _, r := range results {
		base += "name: " + r.Title + "\n"
		base += "endpoint: " + r.Endpoint + "\n"
		base += "status: 【 " + string(r.Status()) + " 】\n"
		base += "check http://172.17.74.34:3000/ for more info\n"
	}

	return base
}
