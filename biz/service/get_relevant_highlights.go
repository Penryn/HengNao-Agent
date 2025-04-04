package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/google/uuid"
	"meeting_agent/biz/dal/mysql"
	"meeting_agent/biz/model"
	"meeting_agent/biz/utils"
	"meeting_agent/conf"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	hertz_gen "meeting_agent/hertz_gen"
)

type GetRelevantHighlightsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetRelevantHighlightsService(Context context.Context, RequestContext *app.RequestContext) *GetRelevantHighlightsService {
	return &GetRelevantHighlightsService{RequestContext: RequestContext, Context: Context}
}

type GetRelevantHighlightsInputs struct {
	UserFavor   []string      `json:"userFavor"`
	EventInform []EventInform `json:"eventInform"`
}

type EventInform struct {
	Name       string   `json:"name"`
	Keywords   []string `json:"keywords"`
	Highlights []string `json:"highlights"`
}

type InnerAgendas struct {
	TopRecommendedAgendas []string `json:"top_recommended_agendas"`
}

func (h *GetRelevantHighlightsService) Run(req *hertz_gen.GetRelevantHighlightsReq) (resp *hertz_gen.GetRelevantHighlightsResp, err error) {
	q := model.NewMeetingQuery(h.Context, mysql.DB)
	client := utils.RestyClient
	u := uuid.New()
	key := conf.GetConf().Api.Key
	secret := conf.GetConf().Api.Secret
	meetings, num, err := q.GetAll(0, 0)
	eventInform := make([]EventInform, 0, num)
	if err != nil {
		return nil, err
	}
	for _, m := range meetings {
		eventInform = append(eventInform, EventInform{
			Name:       m.Name,
			Keywords:   strings.Split(m.KeyWords, "+"),
			Highlights: strings.Split(m.Highlights, "+"),
		})
	}
	inputs := GetRelevantHighlightsInputs{
		UserFavor:   req.UserFavor,
		EventInform: eventInform,
	}
	data := utils.AgentReq{
		SID:    u.String(),
		ID:     conf.GetConf().Api.GetRelevantHighlights,
		Stream: false,
		Inputs: inputs,
	}
	result := &utils.AgentResp{}
	_, err = client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"appKey":       key,
			"sign":         utils.GetSign(key, secret),
		}).
		SetBody(data).
		SetResult(result).
		Post(conf.GetConf().Api.Url)
	var inner InnerAgendas
	err = json.Unmarshal([]byte(result.Data.Results.Output), &inner)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}
	return &hertz_gen.GetRelevantHighlightsResp{
		TopRecommendations: inner.TopRecommendedAgendas,
	}, err
}
