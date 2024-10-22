package packets

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

const (
	TypeRequestPacket           = 0
	TypeResponsePacket          = 1
	TypeTaskCommunicationPacket = 2
)

type BasePacket struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
}

type RequestPacket struct {
	BasePacket
	RequestType string      `json:"requestType"`
	Body        interface{} `json:"body"`
}

const (
	TypeCommandExecution = "CommandExecution"
	TypePlayAudio        = "PlayAudio"
	TypePlayVideo        = "PlayVideo"
	TypePingRequest      = "Ping"
	TypeStreamingRequest = "Streaming"
)

type PingRequest struct{}

type HandShake struct {
	Name         string `json:"name"`
	Os           string `json:"os"`
	Architecture string `json:"architecture"`
}

type CommandExecutionRequest struct {
	Command string `json:"command"`
}

type StreamingRequest struct {
	ChannelId string `json:"channelId"`
	Fps       int    `json:"fps"`
}

type PlayAudioRequest struct {
	URL    string `json:"url"`
	Volume int    `json:"volume"`
}

type PlayVideoRequest struct {
	URL        string `json:"url"`
	Volume     int    `json:"volume"`
	FullScreen bool   `json:"fullScreen"`
}

type ResponsePacket struct {
	BasePacket
	RequestId    string
	ResponseType string      `json:"responseType"`
	Body         interface{} `json:"body"`
}

const (
	TypeTextResponse       = "TextResponse"
	TypeErrorResponse      = "Error"
	TypeTaskLaunchResponse = "TaskLaunch"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type TextResponse struct {
	Text string `json:"text"`
}

type TaskLaunchResponse struct {
	TaskId   string `json:"taskId"`
	TaskType string `json:"taskType"`
}

type TaskCommunicationPackt struct {
	BasePacket
	TaskId    string      `json:"taskId"`
	TaskType  string      `json:"taskType"`
	EventType string      `json:"eventType"`
	Body      interface{} `json:"body"`
}

func ParsePacket(content string) (interface{}, error) {
	contentBytes := []byte(content)
	var packet BasePacket
	if err := json.Unmarshal(contentBytes, &packet); err != nil {
		return nil, err
	}
	switch packet.Type {
	case TypeRequestPacket:
		return ParseRequestPacket(contentBytes)
	case TypeResponsePacket:
		return ParseResponsePacket(contentBytes)
	case TypeTaskCommunicationPacket:
		return ParseTaskCommunicationPacket(contentBytes)
	default:
		return nil, errors.New("invalid json")
	}
}

func ParseRequestPacket(content []byte) (*RequestPacket, error) {
	var packet RequestPacket
	if err := json.Unmarshal(content, &packet); err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(packet.Body)
	if err != nil {
		return nil, err
	}

	switch packet.RequestType {
	case TypeCommandExecution:
		var body CommandExecutionRequest
		if err := json.Unmarshal(bytes, &body); err != nil {
			return nil, err
		}
		packet.Body = body
	case TypePlayAudio:
		var body PlayAudioRequest
		if err := json.Unmarshal(bytes, &body); err != nil {
			return nil, err
		}
		packet.Body = body
	case TypePlayVideo:
		var body PlayVideoRequest
		if err := json.Unmarshal(bytes, &body); err != nil {
			return nil, err
		}
		packet.Body = body
	default:
		return nil, errors.New("invalid json invalid RequestType provided")
	}

	return &packet, nil
}

func ParseResponsePacket(content []byte) (*ResponsePacket, error) {
	var packet ResponsePacket
	if err := json.Unmarshal(content, &packet); err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(packet.Body)
	if err != nil {
		return nil, errors.New("unexpeced internal erorr while rearshaling json")
	}

	switch packet.ResponseType {
	case TypeTextResponse:
		var body TextResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			return nil, err
		}
		packet.Body = body
	case TypeErrorResponse:
		var body ErrorResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			return nil, err
		}
		packet.Body = body
	case TypeTaskLaunchResponse:
		var body TaskLaunchResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			return nil, err
		}
		packet.Body = body
	default:
		return nil, errors.New("invalid json invalid ResponseType proivded")
	}

	return &packet, nil
}

func ParseTaskCommunicationPacket(content []byte) (*TaskCommunicationPackt, error) {
	var packet TaskCommunicationPackt
	if err := json.Unmarshal(content, &packet); err != nil {
		return nil, err
	}

	switch packet.TaskType {
	case TypeCommandExecution:
		err := ParseCommandCommunication(&packet)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid task type")
	}
	return &packet, nil
}

func NewCommunicationPacket(body interface{}, taskId string) *TaskCommunicationPackt {
	packet := TaskCommunicationPackt{
		BasePacket: BasePacket{
			Id:   uuid.NewString(),
			Type: TypeTaskCommunicationPacket,
		},
		TaskId: taskId,
		Body:   body,
	}

	switch body.(type) {
	case CommandInput:
		packet.TaskType = TypeCommandExecution
		packet.EventType = TypeCommandInput
	case CommandOutput:
		packet.TaskType = TypeCommandExecution
		packet.EventType = TypeCommandOutput
	case CommandTerminate:
		packet.TaskType = TypeCommandExecution
		packet.EventType = TypeCommandTerminate
	case AudioContinueRequest:
		packet.TaskType = TypePlayAudio
		packet.EventType = TypeAudioContinueRequest
	case AudioPauseRequest:
		packet.TaskType = TypePlayAudio
		packet.EventType = TypeAudionPauseRequest
	case VideoContinueRequest:
		packet.TaskType = TypePlayVideo
		packet.EventType = TypeVideoContinueRequest
	case VideoPauseRequest:
		packet.TaskType = TypePlayVideo
		packet.EventType = TypeVideoPauseRequest
	case StreamingResumeRequest:
		packet.TaskType = TypeStreamingRequest
		packet.EventType = TypeStreamingResumeRequest
	case StreamingEndRequest:
		packet.TaskType = TypeStreamingRequest
		packet.EventType = TypeStreamingEndRequest
	case StreamingPauseRequest:
		packet.TaskType = TypeStreamingRequest
		packet.EventType = TypeStreamingPauseRequest
	case StreamingChangeFpsRequest:
		packet.TaskType = TypeStreamingRequest
		packet.EventType = TypeStreamingChangeFpsRequest
	default:
		panic("Invalid argument supplied")
	}
	return &packet
}

func NewRequestPacket(request interface{}) *RequestPacket {
	packet := RequestPacket{
		BasePacket: BasePacket{
			Id:   uuid.NewString(),
			Type: TypeRequestPacket,
		},
		Body: request,
	}

	switch request.(type) {
	case PlayAudioRequest:
		packet.RequestType = TypePlayAudio
	case CommandExecutionRequest:
		packet.RequestType = TypeCommandExecution
	case PlayVideoRequest:
		packet.RequestType = TypePlayVideo
	case PingRequest:
		packet.RequestType = TypePingRequest
	case StreamingRequest:
		packet.RequestType = TypeStreamingRequest
	default:
		panic("Invalid argument sppulied")
	}

	return &packet
}

func NewResponsePacket(response interface{}, requestId string) *ResponsePacket {
	packet := ResponsePacket{
		BasePacket: BasePacket{
			Id:   uuid.NewString(),
			Type: TypeResponsePacket,
		},
		RequestId: requestId,
		Body:      response,
	}

	switch response.(type) {
	case TaskLaunchResponse:
		packet.ResponseType = TypeTaskLaunchResponse
	case TextResponse:
		packet.ResponseType = TypeTextResponse
	case ErrorResponse:
		packet.ResponseType = TypeErrorResponse
	default:
		panic("Invalid argument supplied")
	}

	return &packet
}
