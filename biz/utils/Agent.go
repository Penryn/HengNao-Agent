package utils

type AgentReq struct {
	SID    string `json:"sid"`
	ID     string `json:"id"`
	Inputs Inputs `json:"inputs"`
}

type AgentResp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Flag int    `json:"flag"`
	Data Data   `json:"data"`
	TID  string `json:"tid"`
}

type Inputs struct {
	UserFavor   []string      `json:"userFavor"`
	EventInform []EventInform `json:"eventInform"`
}

type EventInform struct {
	Name       string   `json:"name"`
	Keywords   []string `json:"keywords"`
	Highlights []string `json:"highlights"`
}

type Data struct {
	TokenUsage TokenUsage `json:"token_usage"`
	Session    Session    `json:"session"`
	Results    Results    `json:"results"`
}

type TokenUsage struct {
	CompletionTokens         int                        `json:"completion_tokens"`
	PromptTokens             int                        `json:"prompt_tokens"`
	TotalTokens              int                        `json:"total_tokens"`
	CompletionCTokensDetails []CompletionCTokensDetails `json:"completion_tokens_details"`
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
	TopRecommendedAgendas string `json:"top_recommended_agendas"`
}

type InnerAgendas struct {
	TopRecommendedAgendas []string `json:"top_recommended_agendas"`
}
