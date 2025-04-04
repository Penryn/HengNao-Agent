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

type GetMeetingInfoService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetMeetingInfoService(Context context.Context, RequestContext *app.RequestContext) *GetMeetingInfoService {
	return &GetMeetingInfoService{RequestContext: RequestContext, Context: Context}
}

func (h *GetMeetingInfoService) Run(req *hertz_gen.GetMeetingInfoReq) (resp *hertz_gen.GetMeetingInfoResp, err error) {
	q := model.NewMeetingQuery(h.Context, mysql.DB)
	m, err := q.GetById(req.Id)
	if err != nil {
		return nil, err
	}
	return &hertz_gen.GetMeetingInfoResp{
		Info: &hertz_gen.MeetingInfo{
			Id:         m.ID,
			Name:       m.Name,
			Location:   m.Location,
			Time:       m.Time.Format(time.DateTime),
			KeyWords:   strings.Split(m.KeyWords, "+"),
			Highlights: strings.Split(m.Highlights, "+"),
			Content:    m.Content,
			Minutes:    m.Minutes,
		},
	}, nil
}
