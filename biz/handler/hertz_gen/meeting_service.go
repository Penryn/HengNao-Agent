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
