package packets

import (
	"encoding/json"
	"errors"
)

const (
	TypeCommandInput     = "CommandInput"
	TypeCommandOutput    = "CommandOutput"
	TypeCommandTerminate = "CommandTerminate"
	TypeCommandExited    = "CommandExited"
)

type CommandInput struct {
	Input string `json:"input"`
}

type CommandTerminate struct {
}

type CommandOutput struct {
	Output string `json:"output"`
}

type CommandExited struct {
	ExitCode int `json:"exitcode"`
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
		packet.Body = input
	case TypeCommandOutput:
		var input CommandOutput
		if err = json.Unmarshal(bytes, &input); err != nil {
			return err
		}
		packet.Body = input
	case TypeCommandTerminate:
		var input CommandTerminate
		if err = json.Unmarshal(bytes, &input); err != nil {
			return err
		}
		packet.Body = input
	case TypeCommandExited:
		var input CommandExited
		if err = json.Unmarshal(bytes, &input); err != nil {
			return err
		}
		packet.Body = input
	default:
		return errors.New("invalid event type for command execution task")
	}
	return nil
}
