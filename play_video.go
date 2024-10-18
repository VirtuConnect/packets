package packets

import (
	"encoding/json"
	"errors"
)

const (
	TypeVideoPauseRequest    = "VideoPause"
	TypeVideoContinueRequest = "VideoContinue"
	TypeStatusRequest        = "StatusRequest"
)

type VideoPauseRequest struct{}
type VideoContinueRequest struct{}

type StatusRequest struct {
	Code int `json:"code"`
}

func ParsePlayVideoCommunication(packet *TaskCommunicationPackt) error {
	bytes, err := json.Marshal(packet.Body)
	if err != nil {
		return errors.New("internal json parsing error")
	}

	switch packet.EventType {
	case TypeVideoPauseRequest:
		packet.Body = VideoPauseRequest{}
	case TypeVideoContinueRequest:
		packet.Body = VideoContinueRequest{}
	case TypeStatusRequest:
		var req StatusRequest
		if err := json.Unmarshal(bytes, &req); err != nil {
			return err
		}
		packet.Body = req
	default:
		return errors.New("invalid event type provided")
	}
	return nil
}
