package vo

type WorkflowTaskStatus string

const (
	WorkflowTaskStatusPending   = WorkflowTaskStatus("pending")
	WorkflowTaskStatusRunning   = WorkflowTaskStatus("running")
	WorkflowTaskStatusCompleted = WorkflowTaskStatus("completed")
	WorkflowTaskStatusFailed    = WorkflowTaskStatus("failed")
	WorkflowTaskStatusCanceled  = WorkflowTaskStatus("canceled")
)

type WorkflowTask struct {
	Id           uint               `json:"id"`
	OpenID       string             `json:"openid"`
	TaskId       string             `json:"task_id"`
	WorkflowId   string             `json:"workflow_id"`
	WorkflowName string             `json:"workflow_name,omitempty"`
	Score        int                `json:"score"`
	Params       map[string]any     `json:"params"`
	Output       map[string]any     `json:"output"`
	Status       WorkflowTaskStatus `json:"status"`
	Progress     int                `json:"progress"`
	Error        string             `json:"error,omitempty"`
	CreatedAt    int64              `json:"created_at"`
	UpdatedAt    int64              `json:"updated_at"`
}
