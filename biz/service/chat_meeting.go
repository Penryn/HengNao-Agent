package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"meeting_agent/biz/dal/mysql"
	"meeting_agent/biz/model"
	"meeting_agent/biz/utils"
	"meeting_agent/conf"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/google/uuid"
	"github.com/hertz-contrib/websocket"

	hertz_gen "meeting_agent/hertz_gen"

	"github.com/cloudwego/hertz/pkg/app"
)

type ChatMeetingService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewChatMeetingService(Context context.Context, RequestContext *app.RequestContext) *ChatMeetingService {
	return &ChatMeetingService{RequestContext: RequestContext, Context: Context}
}

type chatMeetingInputs struct {
	Input    string    `json:"input"`
	Meetings []meeting `json:"meetings"`
}

type meeting struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	Location   string   `json:"location"`
	Time       string   `json:"time"`
	Keywords   []string `json:"keywords"`
	Highlights []string `json:"highlights"`
}

var upgrader = websocket.HertzUpgrader{}

func (h *ChatMeetingService) Run(req *hertz_gen.ChatMeetingReq) (resp *hertz_gen.ChatMeetingResp, err error) {
	err = upgrader.Upgrade(h.RequestContext, func(conn *websocket.Conn) {
		for {
			// 读取客户端发送的信息
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read: ", err)
				break
			}
			q := model.NewMeetingQuery(h.Context, mysql.DB)
			meetings, num, err := q.GetAll(0, 0)
			client := utils.RestyClient
			u := uuid.New()
			key := conf.GetConf().Api.Key
			secret := conf.GetConf().Api.Secret

			ms := make([]meeting, 0, num)
			for _, m := range meetings {
				ms = append(ms, meeting{
					ID:         m.ID,
					Name:       m.Name,
					Location:   m.Location,
					Time:       m.Time.Format(time.DateTime),
					Keywords:   strings.Split(m.KeyWords, "+"),
					Highlights: strings.Split(m.Highlights, "+"),
				})
			}

			inputs := chatMeetingInputs{
				Input:    string(message),
				Meetings: ms,
			}
			data := utils.AgentReq{
				SID:    u.String(),
				ID:     conf.GetConf().Api.ChatMeeting,
				Stream: true,
				Inputs: inputs,
			}

			// 使用流式处理
			resp, err := client.R().
				SetHeaders(map[string]string{
					"Content-Type": "application/json",
					"appKey":       key,
					"sign":         utils.GetSign(key, secret),
				}).
				SetBody(data).
				SetDoNotParseResponse(true). // 重要：不解析响应
				Post(conf.GetConf().Api.Url)

			if err != nil {
				fmt.Printf("调用会议纪要API出错: %v\n", err)
				return
			}

			// 获取原始响应体
			rawResponse := resp.RawResponse
			defer rawResponse.Body.Close()

			// 创建一个实时读取流
			reader := bufio.NewReader(rawResponse.Body)

			// 实时处理每一行响应
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Printf("读取响应流出错: %v\n", err)
					break
				}

				// 处理 SSE 格式的数据行
				if strings.HasPrefix(line, "data:") {
					jsonStr := strings.TrimPrefix(line, "data:")
					var response utils.AgentFlowResp
					if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
						fmt.Printf("解析JSON错误: %v, line: %s\n", err, line)
						continue
					}

					// 实时发送到 WebSocket
					if err = conn.WriteMessage(mt, []byte(response.Data.Content)); err != nil {
						log.Println("write: ", err)
						break
					}
				}
			}
		}
	})
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	return &hertz_gen.ChatMeetingResp{}, nil
}
