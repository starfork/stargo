package task

import "encoding/json"

// Handler is the function signature for task execution.
type Handler func(*Task) error

func UnmarshalJson(j string) (*Task, error) {
	var task Task
	err := json.Unmarshal([]byte(j), &task)
	return &task, err
}
