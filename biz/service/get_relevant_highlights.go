package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/google/uuid"
	"meeting_agent/biz/utils"
	"meeting_agent/conf"

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

func (h *GetRelevantHighlightsService) Run(req *hertz_gen.GetRelevantHighlightsReq) (resp *hertz_gen.GetRelevantHighlightsResp, err error) {
	client := utils.RestyClient
	u := uuid.New()
	eventInform := make([]utils.EventInform, 0, len(req.EventInform))
	key := conf.GetConf().Api.Key
	secret := conf.GetConf().Api.Secret
	for _, event := range req.EventInform {
		eventInform = append(eventInform, utils.EventInform{
			Name:       event.Name,
			Keywords:   event.Keywords,
			Highlights: event.Highlights,
		})
	}
	inputs := utils.Inputs{
		UserFavor:   req.UserFavor,
		EventInform: eventInform,
	}
	data := utils.AgentReq{
		SID:    u.String(),
		ID:     conf.GetConf().Api.GetRelevantHighlights,
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
	var inner utils.InnerAgendas
	fmt.Println(result.Data.Results.TopRecommendedAgendas)
	err = json.Unmarshal([]byte(result.Data.Results.TopRecommendedAgendas), &inner)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}
	return &hertz_gen.GetRelevantHighlightsResp{
		TopRecommendations: inner.TopRecommendedAgendas,
	}, err
}
