package domain

import "encoding/json"

type MessageEnvelope struct {
	Metadata struct {
		Channel     string `json:"channel"`
		MessageNum  int    `json:"messageNumber"`
		MessageTime string `json:"messageTime"`
		MessageType string `json:"messageType"`
	} `json:"metadata"`
	Message json.RawMessage `json:"message"`
}

const (
	TypeLaunched       = "RocketLaunched"
	TypeSpeedIncreased = "RocketSpeedIncreased"
	TypeSpeedDecreased = "RocketSpeedDecreased"
	TypeExploded       = "RocketExploded"
	TypeMissionChanged = "RocketMissionChanged"
)

type RocketLaunchedPayload struct {
	Type        string `json:"type"`
	LaunchSpeed int64  `json:"launchSpeed"`
	Mission     string `json:"mission"`
}

type RocketSpeedDeltaPayload struct {
	By int64 `json:"by"`
}

type RocketMissionChangedPayload struct {
	NewMission string `json:"newMission"`
}
