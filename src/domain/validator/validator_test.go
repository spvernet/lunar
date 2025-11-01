package validator_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"lunar/src/domain"
	"lunar/src/domain/validator"
)

func makeEnvelope(ch string, num int, when string, kind string) domain.MessageEnvelope {
	var env domain.MessageEnvelope
	env.Metadata.Channel = ch
	env.Metadata.MessageNum = num
	env.Metadata.MessageTime = when
	env.Metadata.MessageType = kind
	return env
}

func TestValidateEnvelope(t *testing.T) {
	v := validator.New()
	now := time.Now().Format(time.RFC3339Nano)

	tests := []struct {
		name string
		env  domain.MessageEnvelope
		want error
	}{
		{
			name: "ok_RFC3339Nano",
			env:  makeEnvelope("channel-1", 1, now, domain.TypeLaunched),
			want: nil,
		},
		{
			name: "ok_RFC3339",
			env:  makeEnvelope("channel-1", 1, time.Now().Format(time.RFC3339), domain.TypeLaunched),
			want: nil,
		},
		{
			name: "invalid_metadata_empty_channel",
			env:  makeEnvelope("", 1, now, domain.TypeLaunched),
			want: validator.ErrInvalidMetadata,
		},
		{
			name: "invalid_metadata_nonpositive_msgnum",
			env:  makeEnvelope("channel-1", 0, now, domain.TypeLaunched),
			want: validator.ErrInvalidMetadata,
		},
		{
			name: "invalid_time",
			env:  makeEnvelope("channel-1", 1, "not-a-time", domain.TypeLaunched),
			want: validator.ErrInvalidMessageTime,
		},
		{
			name: "unknown_type",
			env:  makeEnvelope("channel-1", 1, now, "UnknownType"),
			want: validator.ErrUnknownType,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := v.ValidateEnvelope(tc.env)
			if tc.want == nil && err != nil {
				t.Fatalf("got err=%v, want nil", err)
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("got err=%v, want %v", err, tc.want)
			}
		})
	}
}

func TestValidatePayload(t *testing.T) {
	v := validator.New()

	// Helpers para raw payloads
	raw := func(v any) json.RawMessage {
		b, _ := json.Marshal(v)
		return json.RawMessage(b)
	}

	tests := []struct {
		name string
		kind string
		raw  json.RawMessage
		want error
	}{
		// RocketLaunched
		{
			name: "launched_ok",
			kind: domain.TypeLaunched,
			raw:  raw(domain.RocketLaunchedPayload{Type: "Falcon-9", LaunchSpeed: 100, Mission: "M1"}),
			want: nil,
		},
		{
			name: "launched_invalid_missing_type",
			kind: domain.TypeLaunched,
			raw:  raw(domain.RocketLaunchedPayload{Type: "", LaunchSpeed: 100, Mission: "M1"}),
			want: validator.ErrInvalidPayload,
		},
		{
			name: "launched_invalid_negative_speed",
			kind: domain.TypeLaunched,
			raw:  raw(domain.RocketLaunchedPayload{Type: "F9", LaunchSpeed: -1, Mission: "M1"}),
			want: validator.ErrInvalidPayload,
		},

		// SpeedIncreased
		{
			name: "speed_inc_ok",
			kind: domain.TypeSpeedIncreased,
			raw:  raw(domain.RocketSpeedDeltaPayload{By: 5}),
			want: nil,
		},
		{
			name: "speed_inc_invalid_zero",
			kind: domain.TypeSpeedIncreased,
			raw:  raw(domain.RocketSpeedDeltaPayload{By: 0}),
			want: validator.ErrInvalidPayload,
		},

		// SpeedDecreased
		{
			name: "speed_dec_ok",
			kind: domain.TypeSpeedDecreased,
			raw:  raw(domain.RocketSpeedDeltaPayload{By: 5}),
			want: nil,
		},
		{
			name: "speed_dec_invalid_negative_or_zero",
			kind: domain.TypeSpeedDecreased,
			raw:  raw(domain.RocketSpeedDeltaPayload{By: 0}),
			want: validator.ErrInvalidPayload,
		},

		// MissionChanged
		{
			name: "mission_changed_ok",
			kind: domain.TypeMissionChanged,
			raw:  raw(domain.RocketMissionChangedPayload{NewMission: "M2"}),
			want: nil,
		},
		{
			name: "mission_changed_invalid_empty",
			kind: domain.TypeMissionChanged,
			raw:  raw(domain.RocketMissionChangedPayload{NewMission: ""}),
			want: validator.ErrInvalidPayload,
		},

		// Exploded (no validación extra)
		{
			name: "exploded_ok_empty_payload",
			kind: domain.TypeExploded,
			raw:  raw(struct{}{}),
			want: nil,
		},

		// Unknown type
		{
			name: "unknown_type",
			kind: "UnknownType",
			raw:  raw(struct{}{}),
			want: validator.ErrUnknownType,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := v.ValidatePayload(tc.kind, tc.raw)
			if tc.want == nil && err != nil {
				t.Fatalf("got err=%v, want nil", err)
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("got err=%v, want %v", err, tc.want)
			}
		})
	}
}
