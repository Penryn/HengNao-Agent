package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"meeting_agent/biz/dal/mysql"
	"meeting_agent/biz/dal/redis"
	"meeting_agent/biz/model"
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
	q := model.NewMeetingQuery(h.Context, mysql.DB)
	_, err = q.GetById(req.Id)
	if err != nil {
		return nil, err
	}
	if isProcess(req.Id) {
		return nil, errors.New("会议纪要正在生成中")
	}
	addProcess(req.Id)
	_, err = q.Update(req.Id, model.Meeting{
		Content: req.Content,
		Minutes: " ",
	})
	if err != nil {
		return nil, err
	}

	// 异步调用会议纪要生成
	go func() {
		client := utils.RestyClient
		u := uuid.New()
		key := conf.GetConf().Api.Key
		secret := conf.GetConf().Api.Secret
		inputs := GetMeetingMinutesInputs{
			Input: req.Content,
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
			fmt.Printf("调用会议纪要API出错: %v\n", err)
			return
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
			return
		}

		// 创建新的上下文以便在goroutine中使用
		bgCtx := context.Background()
		newQuery := model.NewMeetingQuery(bgCtx, mysql.DB)

		// 将会议纪要内容存储到数据库
		_, updateErr := newQuery.Update(req.Id, model.Meeting{
			Minutes: minutes,
		})
		if updateErr != nil {
			fmt.Printf("更新会议纪要出错: %v\n", updateErr)
		}
		removeProcess(req.Id)
	}()

	// 立即返回成功响应
	return resp, nil
}

func isProcess(mid uint64) bool {
	q := redis.RedisClient.SIsMember(context.Background(), "meeting_processing", fmt.Sprintf("%d", mid)).Val()
	return q
}

func addProcess(mid uint64) {
	redis.RedisClient.SAdd(context.Background(), "meeting_processing", fmt.Sprintf("%d", mid))
}

func removeProcess(mid uint64) {
	redis.RedisClient.SRem(context.Background(), "meeting_processing", fmt.Sprintf("%d", mid))
}
