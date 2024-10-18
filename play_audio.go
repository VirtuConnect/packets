package packets

import (
	"encoding/json"
	"errors"
)

type AudioPauseRequest struct{}
type AudioContinueRequest struct{}

const (
	TypeAudionPauseRequest   = "AudioPause"
	TypeAudioContinueRequest = "AudioContinue"
)

func ParsePlayAudioCommunication(packet *TaskCommunicationPackt) error {
	bytes, err := json.Marshal(packet.Body)
	if err != nil {
		return errors.New("internal json parsing error")
	}

	switch packet.EventType {
	case TypeAudionPauseRequest:
		packet.Body = AudioPauseRequest{}
	case TypeAudioContinueRequest:
		packet.Body = AudioContinueRequest{}
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
