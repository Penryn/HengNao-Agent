package hertz_gen

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"meeting_agent/biz/service"
	"meeting_agent/biz/utils"
	hertz_gen "meeting_agent/hertz_gen"
)

// GetRelevantHighlights .
// @router /api/meeting/recommendation [POST]
func GetRelevantHighlights(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertz_gen.GetRelevantHighlightsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &hertz_gen.GetRelevantHighlightsResp{}
	resp, err = service.NewGetRelevantHighlightsService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetMeetingMinutes .
// @router /api/meeting/minutes [POST]
func GetMeetingMinutes(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertz_gen.GetMeetingMinutesReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewGetMeetingMinutesService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CreateMeeting .
// @router /api/meeting/create [POST]
func CreateMeeting(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertz_gen.CreateMeetingReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewCreateMeetingService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetMeetingList .
// @router /api/meeting/list [POST]
func GetMeetingList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertz_gen.GetMeetingListReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewGetMeetingListService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetMeetingInfo .
// @router /api/meeting/info [POST]
func GetMeetingInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertz_gen.GetMeetingInfoReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewGetMeetingInfoService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ChatMeeting .
// @router /api/meeting/chat [POST]
func ChatMeeting(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertz_gen.ChatMeetingReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	_, err = service.NewChatMeetingService(ctx, c).Run(&req)

}

// TranslateText .
// @router /api/meeting/translate [GET]
func TranslateText(ctx context.Context, c *app.RequestContext) {
	var err error
	var req hertz_gen.TranslateTextReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewTranslateTextService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
