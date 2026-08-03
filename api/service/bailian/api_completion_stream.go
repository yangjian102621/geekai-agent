package bailian

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// 请求体类型
type CompletionRequest struct {
	Input      CompletionInput      `json:"input"`
	Parameters CompletionParameters `json:"parameters"`
	Debug      interface{}          `json:"debug,omitempty"`
}

type CompletionInput struct {
	Prompt    string      `json:"prompt"`
	SessionId string      `json:"session_id,omitempty"`
	FileList  []FileInfo  `json:"file_list,omitempty"`
	ImageList []ImageInfo `json:"image_list,omitempty"`
}

type FileInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ImageInfo struct {
	Image string `json:"image"`
}

type CompletionParameters struct {
	IncrementalOutput bool `json:"incremental_output"`
	HasThoughts       bool `json:"has_thoughts"`
}

// 响应体类型
type CompletionResponse struct {
	Output    CompletionOutput `json:"output"`
	Usage     CompletionUsage  `json:"usage"`
	RequestId string           `json:"request_id"`
}

type CompletionOutput struct {
	Text           string         `json:"text"`
	FinishReason   string         `json:"finish_reason"`
	SessionId      string         `json:"session_id"`
	Thoughts       []Thought      `json:"thoughts,omitempty"`
	DocReferences  []DocReference `json:"doc_references,omitempty"`
	WorkflowMessage interface{}   `json:"workflow_message,omitempty"`
}

type Thought struct {
	ThoughtName string `json:"thought_name"`
	ActionName  string `json:"action_name"`
	ActionInput string `json:"action_input"`
	Response    string `json:"response"`
}

type DocReference struct {
	DocName  string `json:"doc_name"`
	DocUrl   string `json:"doc_url"`
	ChunkId  string `json:"chunk_id"`
}

type CompletionUsage struct {
	Models []ModelUsage `json:"models"`
}

type ModelUsage struct {
	ModelId     string `json:"model_id"`
	InputTokens int    `json:"input_tokens"`
	OutputTokens int   `json:"output_tokens"`
}

type StreamChannelResponse struct {
	CompletionResponse
	Err error `json:"-"`
}

func (api *API) CompletionStreamRaw(ctx context.Context, req *CompletionRequest) (*http.Response, error) {
	req.Parameters.IncrementalOutput = true
	req.Parameters.HasThoughts = true

	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, "", req)
	if err != nil {
		return nil, err
	}
	return api.c.sendRequest(httpReq)
}

func (api *API) CompletionStream(ctx context.Context, req *CompletionRequest) (chan StreamChannelResponse, error) {
	httpResp, err := api.CompletionStreamRaw(ctx, req)
	if err != nil {
		return nil, err
	}

	streamChannel := make(chan StreamChannelResponse)
	go api.completionStreamHandle(ctx, httpResp, streamChannel)
	return streamChannel, nil
}

func (api *API) completionStreamHandle(ctx context.Context, resp *http.Response, streamChannel chan StreamChannelResponse) {
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
				streamChannel <- StreamChannelResponse{
					Err: fmt.Errorf("error reading line: %w", err),
				}
				return
			}

			if !bytes.HasPrefix(line, []byte("data:")) {
				continue
			}
			line = bytes.TrimPrefix(line, []byte("data:"))

			var streamResp StreamChannelResponse
			if err = json.Unmarshal(line, &streamResp); err != nil {
				streamChannel <- StreamChannelResponse{
					Err: fmt.Errorf("error unmarshalling event: %w", err),
				}
				return
			}

			// 检查是否有错误
			if streamResp.Output.FinishReason == "error" || streamResp.RequestId == "" {
				streamChannel <- StreamChannelResponse{
					Err: fmt.Errorf("error streaming event: %s", string(line)),
				}
				return
			}

			streamChannel <- streamResp

			// finish_reason 为 stop 表示完成
			if streamResp.Output.FinishReason == "stop" {
				return
			}
		}
	}
}
