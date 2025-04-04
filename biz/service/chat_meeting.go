package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/google/uuid"
	"meeting_agent/biz/dal/mysql"
	"meeting_agent/biz/model"
	"meeting_agent/biz/utils"
	"meeting_agent/conf"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	hertz_gen "meeting_agent/hertz_gen"
)

type ChatMeetingService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewChatMeetingService(Context context.Context, RequestContext *app.RequestContext) *ChatMeetingService {
	return &ChatMeetingService{RequestContext: RequestContext, Context: Context}
}

type chatMeetingInputs struct {
	Input    string    `json:"input"`
	Meetings []meeting `json:"meetings"`
}

type meeting struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	Location   string   `json:"location"`
	Time       string   `json:"time"`
	Keywords   []string `json:"keywords"`
	Highlights []string `json:"highlights"`
}

func (h *ChatMeetingService) Run(req *hertz_gen.ChatMeetingReq) (resp *hertz_gen.ChatMeetingResp, err error) {
	q := model.NewMeetingQuery(h.Context, mysql.DB)
	meetings, num, err := q.GetAll(0, 0)
	client := utils.RestyClient
	u := uuid.New()
	key := conf.GetConf().Api.Key
	secret := conf.GetConf().Api.Secret

	ms := make([]meeting, 0, num)
	for _, m := range meetings {
		ms = append(ms, meeting{
			ID:         m.ID,
			Name:       m.Name,
			Location:   m.Location,
			Time:       m.Time.Format(time.DateTime),
			Keywords:   strings.Split(m.KeyWords, "+"),
			Highlights: strings.Split(m.Highlights, "+"),
		})
	}

	inputs := chatMeetingInputs{
		Input:    req.Input,
		Meetings: ms,
	}
	data := utils.AgentReq{
		SID:    u.String(),
		ID:     conf.GetConf().Api.ChatMeeting,
		Stream: true,
		Inputs: inputs,
	}

	r, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"appKey":       key,
			"sign":         utils.GetSign(key, secret),
		}).
		SetBody(data).
		Post(conf.GetConf().Api.Url)

	if err != nil {
		fmt.Printf("调用会议纪要API出错: %v\n", err)
		return
	}
	// 解析响应并提取content内容
	var message string
	scanner := bufio.NewScanner(strings.NewReader(string(r.Body())))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			jsonStr := strings.TrimPrefix(line, "data:")
			var response utils.AgentFlowResp
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				fmt.Printf("解析JSON错误: %v, line: %s\n", err, line)
				continue
			}

			// 提取content并拼接
			message += response.Data.Content
		}
	}
	return &hertz_gen.ChatMeetingResp{
		Output: message,
	}, err
}
