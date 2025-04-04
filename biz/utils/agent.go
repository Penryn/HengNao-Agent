package utils

type AgentReq struct {
	SID    string `json:"sid"`
	ID     string `json:"id"`
	Stream bool   `json:"stream"`
	Inputs any    `json:"inputs"`
}

type AgentResp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Flag int    `json:"flag"`
	Data Data   `json:"data"`
	TID  string `json:"tid"`
}

type Data struct {
	TokenUsage TokenUsage `json:"token_usage"`
	Session    Session    `json:"session"`
	Results    Results    `json:"results"`
}

type TokenUsage struct {
	CompletionTokens         int                      `json:"completion_tokens"`
	PromptTokens             int                      `json:"prompt_tokens"`
	TotalTokens              int                      `json:"total_tokens"`
	CompletionCTokensDetails CompletionCTokensDetails `json:"completion_tokens_details"`
}

type CompletionCTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

type Session struct {
	ID       string    `json:"id"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Results struct {
	Output string `json:"output"`
}

type AgentFlowResp struct {
	Msg  string   `json:"msg"`
	Code int      `json:"code"`
	Data FlowData `json:"data"`
}

type FlowData struct {
	Model     string         `json:"model"`
	Type      string         `json:"type"`
	From      string         `json:"from"`
	Name      string         `json:"name"`
	TimeStamp string         `json:"time_stamp"`
	MessageID string         `json:"message_id"`
	Content   string         `json:"content"`
	Result    map[string]any `json:"result"`
	NodeID    string         `json:"node_id"`
}
