package models

// ConnectionItem describes a single connection destination in n8n.
type ConnectionItem struct {
	Node  string `json:"node"`
	Type  string `json:"type"`
	Index int    `json:"index"`
}

// WorkflowTag represents classification tags assigned to a workflow.
type WorkflowTag struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// WorkflowSettings holds execution and persistence settings for a workflow.
type WorkflowSettings struct {
	ErrorWorkflow            string                 `json:"errorWorkflow,omitempty"`
	ExecutionTimeout         int                    `json:"executionTimeout,omitempty"`
	SaveExecutionProgress    bool                   `json:"saveExecutionProgress,omitempty"`
	SaveManualExecutions     bool                   `json:"saveManualExecutions,omitempty"`
	SaveDataErrorExecution   string                 `json:"saveDataErrorExecution,omitempty"`
	SaveDataSuccessExecution string                 `json:"saveDataSuccessExecution,omitempty"`
	ExecutionOrder           string                 `json:"executionOrder,omitempty"`
	CallerPolicy             string                 `json:"callerPolicy,omitempty"`
	Timezone                 string                 `json:"timezone,omitempty"`
	Raw                      map[string]interface{} `json:"-"`
}

// Node represents an individual step or integration within the n8n canvas.
type Node struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Type             string                 `json:"type"`
	TypeVersion      float64                `json:"typeVersion"`
	Position         []float64              `json:"position,omitempty"`
	Disabled         bool                   `json:"disabled,omitempty"`
	Notes            string                 `json:"notes,omitempty"`
	NotesInFlow      bool                   `json:"notesInFlow,omitempty"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
	Credentials      map[string]interface{} `json:"credentials,omitempty"`
	RetryOnFail      bool                   `json:"retryOnFail,omitempty"`
	MaxTries         int                    `json:"maxTries,omitempty"`
	WaitBetweenTries int                    `json:"waitBetweenTries,omitempty"`
	ContinueOnFail   bool                   `json:"continueOnFail,omitempty"`
	AlwaysOutputData bool                   `json:"alwaysOutputData,omitempty"`
	ExecuteOnce      bool                   `json:"executeOnce,omitempty"`
	OnError          string                 `json:"onError,omitempty"`
	Raw              map[string]interface{} `json:"-"`
}

// Workflow is the root structure of an exported n8n workflow definition.
type Workflow struct {
	ID          string                                  `json:"id,omitempty"`
	Name        string                                  `json:"name"`
	Nodes       []*Node                                 `json:"nodes"`
	Connections map[string]map[string][][]ConnectionItem `json:"connections"`
	Settings    WorkflowSettings                        `json:"settings"`
	PinData     map[string]interface{}                  `json:"pinData,omitempty"`
	Tags        []WorkflowTag                           `json:"tags,omitempty"`
	VersionID   string                                  `json:"versionId,omitempty"`
	Meta        map[string]interface{}                  `json:"meta,omitempty"`
	FilePath    string                                  `json:"-"`
	RawBytes    []byte                                  `json:"-"`
}

// FindNodeByName searches for a node matching the specified name.
func (w *Workflow) FindNodeByName(name string) *Node {
	for _, n := range w.Nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

// FindNodeByID searches for a node matching the specified unique ID.
func (w *Workflow) FindNodeByID(id string) *Node {
	for _, n := range w.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}
