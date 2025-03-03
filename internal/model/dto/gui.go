package dto

type GUITask struct {
    ID     string      `json:"id"`
    Type   string      `json:"type"`
    Params interface{} `json:"params,omitempty"`
}