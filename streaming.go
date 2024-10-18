package packets

import (
	"encoding/json"
	"errors"
)

const (
	TypeStreamingPauseRequest     = "Pause"
	TypeStreamingResumeRequest    = "Resume"
	TypeStreamingEndRequest       = "End"
	TypeStreamingChangeFpsRequest = "ChangeFps"
)

type StreamingPauseRequest struct{}
type StreamingResumeRequest struct{}
type StreamingEndRequest struct{}
type StreamingChangeFpsRequest struct {
	Fps int `json:"fps"`
}

func ParseStreamingCommunication(packet *TaskCommunicationPackt) error {
	bytes, err := json.Marshal(packet.Body)
	if err != nil {
		return err
	}

	switch packet.EventType {
	case TypeStreamingEndRequest:
		packet.Body = StreamingEndRequest{}
	case TypeStreamingPauseRequest:
		packet.Body = StreamingPauseRequest{}
	case TypeStreamingResumeRequest:
		packet.Body = StreamingResumeRequest{}
	case TypeStreamingChangeFpsRequest:
		var body StreamingChangeFpsRequest
		if err := json.Unmarshal(bytes, &body); err != nil {
			return err
		}
		packet.Body = body
	default:
		return errors.New("invalid event type for streaming")
	}
	return nil
}
