package aliqwen

// 开通地址 https://dashscope.console.aliyun.com/apiKey

const ApiBaseUrl = "https://dashscope.aliyuncs.com"
const ApiVersion = "v1"

const ChatMessageRoleAssistant = "assistant"
const ChatMessageRoleSystem = "system"
const ChatMessageRoleUser = "user"

// request

type RequestBody struct {
	Model      string      `json:"model"`
	Input      Input       `json:"input"`
	Parameters *Parameters `json:"parameters"`
}

type Input struct {
	Messages []*Messages `json:"messages"`
}

type Parameters struct {
	EnableSearch bool `json:"enable_search"`
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// response

type ResponseBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Output    Output `json:"output"`
	Usage     Usage  `json:"usage"`
	RequestID string `json:"request_id"`
}

type Output struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	OutputTokens int `json:"output_tokens"`
	InputTokens  int `json:"input_tokens"`
}
