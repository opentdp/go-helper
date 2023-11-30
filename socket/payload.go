package socket

import (
	"encoding/json"
)

type PlainData struct {
	Method  string
	TaskId  uint
	Success bool
	Message string
	Payload any
}

func (d *PlainData) GetPayload(v any) error {

	payload, err := json.Marshal(d.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(payload, v)

}
