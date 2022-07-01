package message

const (
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
)

type CommonResponse struct {
	Status   string      `json:"status"`
	Messages string      `json:"messages"`
	Result   interface{} `json:"result"`
}

type GetWorkFlowsStatusRequest struct {
	WorkflowName string `json:"work_flow_name"`
}

type GetWorkFlowStatusResponse struct {
	Items []WorkFlowInfo `json:"items"`
}

type WorkFlowInfo struct {
	WorkflowName string `json:"work_flow_name"`
	StatusPhase  string `json:"phase_status"`
	S2ID         string `json:"s2id"`
}
