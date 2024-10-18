package packets

import (
	"encoding/json"
	"errors"
)

const (
	TypeCommandInput     = "CommandInput"
	TypeCommandOutput    = "CommandOutput"
	TypeCommandTerminate = "CommandTerminate"
)

type CommandInput struct {
	Input string `json:"input"`
}

type CommandTerminate struct {
}

type CommandOutput struct {
	Output string `json:"output"`
}

func ParseCommandCommunication(packet *TaskCommunicationPackt) error {
	bytes, err := json.Marshal(packet.Body)
	if err != nil {
		return errors.New("internal json parsing error")
	}

	switch packet.EventType {
	case TypeCommandInput:
		var input CommandInput
		if err = json.Unmarshal(bytes, &input); err != nil {
			return err
		}
		packet.Body = bytes
	case TypeCommandOutput:
		var input CommandOutput
		if err = json.Unmarshal(bytes, &input); err != nil {
			return err
		}
		packet.Body = bytes
	case TypeCommandTerminate:
		var input CommandTerminate
		if err = json.Unmarshal(bytes, &input); err != nil {
			return err
		}
		packet.Body = bytes
	default:
		return errors.New("invalid event type for command execution task")
	}
	return nil
}
