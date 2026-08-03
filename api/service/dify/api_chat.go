package dify

import (
	"context"
	"net/http"
)

type ResponseMode string

const (
	ResponseModeBlocking  ResponseMode = "blocking"
	ResponseModeStreaming ResponseMode = "streaming"
)

type ChatMessageRequest struct {
	Inputs         map[string]interface{} `json:"inputs"`
	Query          string                 `json:"query"`
	ResponseMode   ResponseMode           `json:"response_mode"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	User           string                 `json:"user"`
}

type ChatMessageResponse struct {
	ID             string `json:"id"`
	Answer         string `json:"answer"`
	ConversationID string `json:"conversation_id"`
	CreatedAt      int    `json:"created_at"`
}

/* Create chat message
 * Create a new conversation message or continue an existing dialogue.
 */
func (api *API) ChatMessages(ctx context.Context, req *ChatMessageRequest) (resp *ChatMessageResponse, err error) {
	req.ResponseMode = ResponseModeBlocking

	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, "/v1/chat-messages", req)
	if err != nil {
		return
	}
	err = api.c.sendJSONRequest(httpReq, &resp)
	return
}
