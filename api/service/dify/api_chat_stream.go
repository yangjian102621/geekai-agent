package dify

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type MessageEvent string

const (
	MessageEventMessageDelta     MessageEvent = "message_delta"
	MessageEventMessageCompleted MessageEvent = "message_completed"
	MessageEventAgentMessage     MessageEvent = "agent_message"
	MessageEventMessageEnd       MessageEvent = "message_end"
	MessageEventAgentThought     MessageEvent = "agent_thought"
	MessageEventAgentToolStart   MessageEvent = "agent_tool_start"
	MessageEventAgentToolEnd     MessageEvent = "agent_tool_end"
)

type ChatMessageStreamResponse struct {
	Event          MessageEvent                 `json:"event"`
	TaskID         string                       `json:"task_id"`
	ID             string                       `json:"id"`
	Answer         string                       `json:"answer"`
	CreatedAt      int64                        `json:"created_at"`
	ConversationID string                       `json:"conversation_id"`
	Tool           string                       `json:"tool,omitempty"`
	ToolName       string                       `json:"tool_name,omitempty"`
	Observation    string                       `json:"observation,omitempty"`
	ToolLabels     map[string]map[string]string `json:"tool_labels,omitempty"`
	Thought        string                       `json:"thought,omitempty"`
}

type ChatMessageStreamChannelResponse struct {
	ChatMessageStreamResponse
	Err error `json:"-"`
}

func (api *API) ChatMessagesStreamRaw(ctx context.Context, req *ChatMessageRequest) (*http.Response, error) {
	req.ResponseMode = ResponseModeStreaming

	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, "/v1/chat-messages", req)
	if err != nil {
		return nil, err
	}
	return api.c.sendRequest(httpReq)
}

func (api *API) ChatMessagesStream(ctx context.Context, req *ChatMessageRequest) (chan ChatMessageStreamChannelResponse, error) {
	httpResp, err := api.ChatMessagesStreamRaw(ctx, req)
	if err != nil {
		return nil, err
	}

	streamChannel := make(chan ChatMessageStreamChannelResponse)
	go api.chatMessagesStreamHandle(ctx, httpResp, streamChannel)
	return streamChannel, nil
}

func (api *API) chatMessagesStreamHandle(ctx context.Context, resp *http.Response, streamChannel chan ChatMessageStreamChannelResponse) {
	defer resp.Body.Close()
	defer close(streamChannel)

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				streamChannel <- ChatMessageStreamChannelResponse{
					Err: fmt.Errorf("error reading line: %w", err),
				}
				return
			}

			if !bytes.HasPrefix(line, []byte("data:")) {
				continue
			}
			line = bytes.TrimPrefix(line, []byte("data:"))
			var resp ChatMessageStreamChannelResponse
			if err = json.Unmarshal(line, &resp); err != nil {
				streamChannel <- ChatMessageStreamChannelResponse{
					Err: fmt.Errorf("error unmarshalling event: %w", err),
				}
				return
			} else if resp.Event == "error" {
				streamChannel <- ChatMessageStreamChannelResponse{
					Err: errors.New("error streaming event: " + string(line)),
				}
				return
			} else if resp.Tool != "" {
				// 提取工具的名称
				resp.ToolName = resp.ToolLabels[resp.Tool]["zh_Hans"]
				if resp.Observation == "" {
					resp.Event = MessageEventAgentToolStart
				} else {
					resp.Event = MessageEventAgentToolEnd
				}
			} else if resp.Event == MessageEventAgentMessage {
				resp.Event = MessageEventMessageDelta
			} else if resp.Event == MessageEventMessageEnd {
				return
			}
			streamChannel <- resp
		}
	}
}
