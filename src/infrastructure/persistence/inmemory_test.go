package persistence_test

import (
	"encoding/json"
	"testing"
	"time"

	"lunar/src/domain"
	"lunar/src/infrastructure/persistence"
)

func TestMissionChanged_OlderMessageDoesNotOverrideNewer(t *testing.T) {
	store := persistence.NewMemoryStore()
	ch := "193270a9-c9cf-404a-8f83-838e71d9ae67"

	// 1) Primero llega una MissionChanged más NUEVA (num=5)
	envNew := makeEnv(ch, 5, "2022-02-02T19:39:10.000000+01:00",
		domain.TypeMissionChanged, domain.RocketMissionChangedPayload{
			NewMission: "MISSION_NEW",
		})
	if err := store.Apply(envNew); err != nil {
		t.Fatalf("apply newer mission failed: %v", err)
	}

	// 2) Luego llega una MissionChanged ANTIGUA (num=3) → NO debe pisar
	envOld := makeEnv(ch, 3, "2022-02-02T19:39:08.000000+01:00",
		domain.TypeMissionChanged, domain.RocketMissionChangedPayload{
			NewMission: "MISSION_OLD",
		})
	if err := store.Apply(envOld); err != nil {
		t.Fatalf("apply older mission failed: %v", err)
	}

	// 3) Verificamos que mantiene la misión nueva y el last msg num
	got, ok, err := store.Get(ch)
	if err != nil {
		t.Fatalf("get rocket failed: %v", err)
	}
	if !ok {
		t.Fatalf("rocket not found")
	}
	if got.Mission != "MISSION_NEW" {
		t.Errorf("mission was overwritten by older message; got=%s want=MISSION_NEW", got.Mission)
	}
	if got.LastMsgNum != 5 {
		t.Errorf("LastMsgNum incorrect; got=%d want=5", got.LastMsgNum)
	}
}

func TestIdempotency_DuplicateMessageIgnored(t *testing.T) {
	store := persistence.NewMemoryStore()
	ch := "193270a9-c9cf-404a-8f83-838e71d9ae67"

	// Mismo mensaje (num=2) dos veces
	env := makeEnv(ch, 2, "2022-02-02T19:39:06.000000+01:00",
		domain.TypeSpeedIncreased, domain.RocketSpeedDeltaPayload{By: 300})
	if err := store.Apply(env); err != nil {
		t.Fatalf("apply first failed: %v", err)
	}
	if err := store.Apply(env); err != nil { // duplicado
		t.Fatalf("apply duplicate failed: %v", err)
	}

	got, ok, err := store.Get(ch)
	if err != nil || !ok {
		t.Fatalf("get rocket failed: %v ok=%v", err, ok)
	}
	// Debe sumar una sola vez
	if got.Speed != 300 {
		t.Errorf("idempotency fail; speed=%d want=300", got.Speed)
	}
}

func TestDeltas_CommutativeWithOutOfOrder(t *testing.T) {
	store := persistence.NewMemoryStore()
	ch := "193270a9-c9cf-404a-8f83-838e71d9ae67"

	// Llegan desordenados: +500 (num=4) y -200 (num=1)
	envPlus := makeEnv(ch, 4, time.Now().Format(time.RFC3339Nano),
		domain.TypeSpeedIncreased, domain.RocketSpeedDeltaPayload{By: 500})
	envMinus := makeEnv(ch, 1, time.Now().Add(-time.Second).Format(time.RFC3339Nano),
		domain.TypeSpeedDecreased, domain.RocketSpeedDeltaPayload{By: 200})

	if err := store.Apply(envPlus); err != nil {
		t.Fatalf("apply plus failed: %v", err)
	}
	if err := store.Apply(envMinus); err != nil {
		t.Fatalf("apply minus failed: %v", err)
	}

	got, ok, err := store.Get(ch)
	if err != nil || !ok {
		t.Fatalf("get rocket failed: %v ok=%v", err, ok)
	}
	// Conmutativo: 0 + 500 - 200 = 300, sin importar el orden
	if got.Speed != 300 {
		t.Errorf("commutative deltas fail; speed=%d want=300", got.Speed)
	}
	// LastMsgNum debe ser el mayor (4)
	if got.LastMsgNum != 4 {
		t.Errorf("LastMsgNum incorrect; got=%d want=4", got.LastMsgNum)
	}
}

// ---------- helpers ----------

func makeEnv(channel string, num int, when string, kind string, payload any) domain.MessageEnvelope {
	raw, _ := json.Marshal(payload)
	env := domain.MessageEnvelope{}
	env.Metadata.Channel = channel
	env.Metadata.MessageNum = num
	env.Metadata.MessageTime = when
	env.Metadata.MessageType = kind
	env.Message = raw
	return env
}
