package service

import (
	"context"
	"meeting_agent/biz/dal/mysql"
	"meeting_agent/biz/model"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	hertz_gen "meeting_agent/hertz_gen"
)

type CreateMeetingService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateMeetingService(Context context.Context, RequestContext *app.RequestContext) *CreateMeetingService {
	return &CreateMeetingService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateMeetingService) Run(req *hertz_gen.CreateMeetingReq) (resp *hertz_gen.CreateMeetingResp, err error) {
	q := model.NewMeetingQuery(h.Context, mysql.DB)
	t, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		return nil, err
	}
	m, err := q.Create(model.Meeting{
		Name:       req.Name,
		Location:   req.Location,
		Time:       t,
		KeyWords:   strings.Join(req.KeyWords, "+"),
		Highlights: strings.Join(req.Highlights, "+"),
		Content:    "",
		Minutes:    "",
	})
	if err != nil {
		return nil, err
	}
	return &hertz_gen.CreateMeetingResp{
		Id: m.ID,
	}, err
}
