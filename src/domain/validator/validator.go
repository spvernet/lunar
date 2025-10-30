package validator

import (
	"encoding/json"
	"errors"
	"time"

	"lunar/src/domain"
)

type payloadValidatorFunc func(json.RawMessage) error

var (
	ErrInvalidMetadata    = errors.New("invalid metadata")
	ErrInvalidMessageTime = errors.New("invalid messageTime (RFC3339/RFC3339Nano)")
	ErrUnknownType        = errors.New("unknown messageType")
	ErrInvalidPayload     = errors.New("invalid payload")
)

var defaultPayloadValidators = map[string]payloadValidatorFunc{
	domain.TypeLaunched:       validateLaunched,
	domain.TypeSpeedIncreased: validateSpeedDeltaPositive,
	domain.TypeSpeedDecreased: validateSpeedDeltaPositive,
	domain.TypeMissionChanged: validateMissionChanged,
	domain.TypeExploded:       validateExplodedNoop,
}

type Validator interface {
	ValidateEnvelope(env domain.MessageEnvelope) error
	ValidatePayload(kind string, raw json.RawMessage) error
}
type validator struct{}

func New() Validator { return validator{} }

func (validator) ValidateEnvelope(env domain.MessageEnvelope) error {
	if env.Metadata.Channel == "" || env.Metadata.MessageNum <= 0 || env.Metadata.MessageType == "" {
		return ErrInvalidMetadata
	}
	// Acepta RFC3339 y RFC3339Nano
	if _, err := time.Parse(time.RFC3339Nano, env.Metadata.MessageTime); err != nil {
		if _, err2 := time.Parse(time.RFC3339, env.Metadata.MessageTime); err2 != nil {
			return ErrInvalidMessageTime
		}
	}
	if _, ok := defaultPayloadValidators[env.Metadata.MessageType]; !ok {
		return ErrUnknownType
	}
	return nil
}
func (validator) ValidatePayload(kind string, raw json.RawMessage) error {
	fn, ok := defaultPayloadValidators[kind]
	if !ok {
		return ErrUnknownType
	}
	return fn(raw)

}

func validateLaunched(raw json.RawMessage) error {
	var p domain.RocketLaunchedPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return ErrInvalidPayload
	}
	if p.Type == "" || p.Mission == "" || p.LaunchSpeed < 0 {
		return ErrInvalidPayload
	}
	return nil
}

func validateSpeedDeltaPositive(raw json.RawMessage) error {
	var p domain.RocketSpeedDeltaPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return ErrInvalidPayload
	}
	if p.By <= 0 {
		return ErrInvalidPayload
	}
	return nil
}

func validateMissionChanged(raw json.RawMessage) error {
	var p domain.RocketMissionChangedPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return ErrInvalidPayload
	}
	if p.NewMission == "" {
		return ErrInvalidPayload
	}
	return nil
}

func validateExplodedNoop(_ json.RawMessage) error { return nil }
