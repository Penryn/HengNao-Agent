package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"math"
	"meeting_agent/biz/dal/mysql"
	"meeting_agent/biz/model"
	hertz_gen "meeting_agent/hertz_gen"
	"strings"
)

type GetMeetingListService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetMeetingListService(Context context.Context, RequestContext *app.RequestContext) *GetMeetingListService {
	return &GetMeetingListService{RequestContext: RequestContext, Context: Context}
}

func (h *GetMeetingListService) Run(req *hertz_gen.GetMeetingListReq) (resp *hertz_gen.GetMeetingListResp, err error) {
	q := model.NewMeetingQuery(h.Context, mysql.DB)
	meetings, num, err := q.GetAll(req.PageNum, req.PageSize)
	if err != nil {
		return nil, err
	}

	ms := make([]*hertz_gen.MeetingInfo, len(meetings))
	for i, m := range meetings {
		ms[i] = &hertz_gen.MeetingInfo{
			Id:         m.ID,
			Name:       m.Name,
			Location:   m.Location,
			Time:       m.Time.Format("2006-01-02 15:04:05"),
			KeyWords:   []string{},
			Highlights: []string{},
			Content:    m.Content,
			Minutes:    m.Minutes,
		}
		if len(m.KeyWords) > 0 {
			ms[i].KeyWords = strings.Split(m.KeyWords, "+")
		}
		if len(m.Highlights) > 0 {
			ms[i].Highlights = strings.Split(m.Highlights, "+")
		}
	}
	return &hertz_gen.GetMeetingListResp{
		Total:       uint64(math.Ceil(float64(num) / float64(req.PageSize))),
		MeetingList: ms,
	}, nil
}
