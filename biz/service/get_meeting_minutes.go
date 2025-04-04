package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"meeting_agent/biz/utils"
	"meeting_agent/conf"
	"strings"

	"github.com/google/uuid"

	hertz_gen "meeting_agent/hertz_gen"

	"github.com/cloudwego/hertz/pkg/app"
)

type GetMeetingMinutesService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetMeetingMinutesService(Context context.Context, RequestContext *app.RequestContext) *GetMeetingMinutesService {
	return &GetMeetingMinutesService{RequestContext: RequestContext, Context: Context}
}

type GetMeetingMinutesInputs struct {
	Input string `json:"input"`
}

func (h *GetMeetingMinutesService) Run(req *hertz_gen.GetMeetingMinutesReq) (resp *hertz_gen.GetMeetingMinutesResp, err error) {
	client := utils.RestyClient
	u := uuid.New()
	key := conf.GetConf().Api.Key
	secret := conf.GetConf().Api.Secret
	inputs := GetMeetingMinutesInputs{
		Input: req.Input,
	}
	data := utils.AgentReq{
		SID:    u.String(),
		ID:     conf.GetConf().Api.GetMeetingMinutes,
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
		return nil, err
	}

	// 解析响应并提取content内容
	var minutes string
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
			minutes += response.Data.Content
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("扫描响应时出错: %v\n", err)
	}

	// 返回会议纪要内容
	resp = &hertz_gen.GetMeetingMinutesResp{
		Output: minutes,
	}

	return resp, nil
}
