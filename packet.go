package packets

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Packet struct {
	PacketType string      `json:"packetType"`
	Body       interface{} `json:"body"`
}

type UnparsedPacket struct {
	PacketType string                 `json:"packetType"`
	Body       map[string]interface{} `json:"body"`
}

const (
	TypePacketReadingError    = "RedingError"
	TypePacketParsingError    = "ParsingError"
	TypePacketProcessingError = "ProcessingError"
	TypeCommandResponsePacket = "CommandResponse"
	TypeStatusResponsePacket  = "StatusResponse"
	TypeTaskLaunchPacket      = "TaskLaunchPacket"
	TypeErrorPacket           = "Error"
	TypeRequestResponsePacket = "RequestResponsePacket"
	TypeTaskResponsePacket    = "TaskResponsePacket"
	TypeRequestPacket         = "RequestPacket"
)

var PacketReadingErrorPacket = Packet{
	PacketType: TypePacketReadingError,
	Body:       "Failed to read packet",
}

var PacketParsingErrorPacket = Packet{
	PacketType: TypePacketParsingError,
	Body:       "Failed to parse the packets body (missing fields or invalid data)",
}

var HelloWorldPacket = Packet{
	PacketType: "HelloWorld",
	Body:       "Hello its me",
}

var PacketProcessingErrorPacket = Packet{
	PacketType: TypePacketParsingError,
	Body:       "Failed to process package. Unkown Failure",
}

func ParsePacket(packet *UnparsedPacket) (*Packet, error) {
	result := &Packet{PacketType: packet.PacketType}

	//convert the unparsed body into json
	jsonData, err := json.Marshal(packet.Body)
	if err != nil {
		return nil, err
	}

	//switch the different cases
	switch packet.PacketType {
	case TypeRequestResponsePacket:
		var packet UnparsedRequestResponsePacket
		err = json.Unmarshal(jsonData, &packet)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to ResponsePacket: %w", err)
		}
		body, err := ParseRequestResponsePacket(&packet)
		if err != nil {
			return nil, err
		}
		result.Body = *body
	case TypeTaskResponsePacket:
		var packet UnparsedTaskResponsePacket
		err = json.Unmarshal(jsonData, &packet)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to ResponsePacket: %w", err)
		}
		body, err := ParseTaskResponsePacket(&packet)
		if err != nil {
			return nil, err
		}
		result.Body = *body
	case TypeRequestPacket:
		var packet UnparsedRequestPacket
		err = json.Unmarshal(jsonData, &packet)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to ResponsePacket: %w", err)
		}
		body, err := ParseRequestPacket(&packet)
		if err != nil {
			return nil, err
		}
		result.Body = *body
	case TypeErrorPacket:
		var errorResponse ErrorPacket
		err = json.Unmarshal(jsonData, &errorResponse)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to ErrorPacket: %w", err)
		}
		result.Body = errorResponse
	default:
		return nil, fmt.Errorf("invalid packet content type `%s`", packet.PacketType)
	}
	return result, nil
}

type UnparsedRequestResponsePacket struct {
	Type      string      `json:"type"`
	RequestId string      `json:"requestId"`
	Body      interface{} `json:"body"`
}

type UnparsedTaskResponsePacket struct {
	Type   string      `json:"type"`
	TaskId string      `json:"taskId"`
	Body   interface{} `json:"body"`
}

type UnparsedRequestPacket struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

func ParseRequestPacket(packet *UnparsedRequestPacket) (*RequestPacket, error) {
	jsonData, err := json.Marshal(packet.Body)
	if err != nil {
		return nil, err
	}
	var result = RequestPacket{Type: packet.Type}
	switch packet.Type {
	case TypePlayVideoRequestPacket:
		var packet PlayVideoRequestPacket
		err = json.Unmarshal(jsonData, &packet)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to CommandResponsePacket: %w", err)
		}
		result.Body = packet
	case TypePlayAudioRequestPacket:
		var packet PlayAudioRequestPacket
		err = json.Unmarshal(jsonData, &packet)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to CommandResponsePacket: %w", err)
		}
		result.Body = packet
	case TypePlayYoutubeRequestPacket:
		var packet PlayYoutubeRequestPacket
		err = json.Unmarshal(jsonData, &packet)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to CommandResponsePacket: %w", err)
		}
		result.Body = packet
	case TypeRunCommandRequestPacket:
		var packet RunCommandRequestPacket
		err = json.Unmarshal(jsonData, &packet)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to CommandResponsePacket: %w", err)
		}
		result.Body = packet
	default:
		return nil, fmt.Errorf("invalid packet content type (request) `%s`", packet.Type)
	}
	return &result, nil
}

func ParseRequestResponsePacket(packet *UnparsedRequestResponsePacket) (*RequestResponsePacket, error) {
	//convert the unparsed body into json
	jsonData, err := json.Marshal(packet.Body)
	if err != nil {
		return nil, err
	}

	var result = RequestResponsePacket{RequestId: packet.RequestId, Type: packet.Type}
	parsedBody, err := parseResponses(packet.Type, jsonData)

	if err != nil {
		return nil, err
	}

	result.Body = parsedBody
	return &result, nil
}

func parseResponses(tp string, jsonData []byte) (interface{}, error) {
	switch tp {
	case TypeCommandResponsePacket:
		var commandResponse CommandResponsePacket
		err := json.Unmarshal(jsonData, &commandResponse)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to CommandResponsePacket: %w", err)
		}

		return commandResponse, nil
	case TypeStatusResponsePacket:
		var statusResponse StatusResponsePacket
		err := json.Unmarshal(jsonData, &statusResponse)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to CommandResponsePacket: %w", err)
		}
		return statusResponse, nil
	case TypeTaskLaunchPacket:
		var launchPacket TaskLaunchPacket
		err := json.Unmarshal(jsonData, &launchPacket)

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling body to CommandResponsePacket: %w", err)
		}
		return launchPacket, nil
	default:
		return nil, fmt.Errorf("invalid packet content type (response) `%s`", tp)
	}
}

func ParseTaskResponsePacket(packet *UnparsedTaskResponsePacket) (*TaskResponsePacket, error) {
	//convert the unparsed body into json
	jsonData, err := json.Marshal(packet.Body)
	if err != nil {
		return nil, err
	}

	var result = TaskResponsePacket{TaskId: packet.TaskId, Type: packet.Type}
	parsedBody, err := parseResponses(packet.Type, jsonData)

	if err != nil {
		return nil, err
	}

	result.Body = parsedBody
	return &result, nil
}

func ReadPacket(conn *websocket.Conn) (*Packet, error) {
	var packet UnparsedPacket
	if err := conn.ReadJSON(&packet); err != nil {
		return nil, err
	}
	return ParsePacket(&packet)
}

type Status int8

const (
	StatusSuccess    = Status(0)
	StatusFailure    = Status(1)
	StatusPending    = Status(2)
	StatusNotFound   = Status(3)
	StatusNotAllowed = Status(4)
)

var stringToStatus = map[string]Status{
	"Success":    StatusSuccess,
	"Failure":    StatusFailure,
	"Pending":    StatusPending,
	"NotFound":   StatusNotFound,
	"NotAllowed": StatusNotAllowed,
}

var statusToString = map[Status]string{
	StatusSuccess:    "Success",
	StatusFailure:    "Failure",
	StatusPending:    "Pending",
	StatusNotFound:   "NotFound",
	StatusNotAllowed: "NotAllowed",
}

func ParseStatus(input string) (Status, error) {
	status, ok := stringToStatus[input]
	if !ok {
		return StatusFailure, fmt.Errorf("failed to parse status")
	}
	return status, nil
}

func (s Status) ToString() string {
	return statusToString[s]
}
