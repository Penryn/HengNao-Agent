package service

import (
	"context"
	"meeting_agent/biz/utils"

	"github.com/cloudwego/hertz/pkg/app"
	hertz_gen "meeting_agent/hertz_gen"
)

type TranslateTextService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewTranslateTextService(Context context.Context, RequestContext *app.RequestContext) *TranslateTextService {
	return &TranslateTextService{RequestContext: RequestContext, Context: Context}
}

func (h *TranslateTextService) Run(req *hertz_gen.TranslateTextReq) (resp *hertz_gen.TranslateTextResp, err error) {
	translated, err := utils.Translate(req.Text, req.TargetLanguage)
	if err != nil {
		return nil, err
	}
	return &hertz_gen.TranslateTextResp{
		TranslatedText: translated,
	}, nil
}
