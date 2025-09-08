package malenkyi

import (
	"errors"
	"testing"
	"time"
)

func TestNewGenerator(t *testing.T) {
	if _, err := NewGenerator(time.Now().Add(-time.Hour*24*365), uint16(0)); err != nil {
		t.Fatalf("using a historical epoch should not return an error but got: %v", err)
	}

	if _, err := NewGenerator(time.Now().Add(time.Hour), uint16(0)); !errors.Is(err, ErrInvalidEpoch) {
		t.Fatalf("using a future epoch should've returned error %v but got %v", ErrInvalidEpoch, err)
	}

	if _, err := NewGenerator(time.Now().Add(-time.Hour), uint16(1023)); err != nil {
		t.Fatalf("using machine ID 1023 should not have returned an error but got: %v", err)
	}

	if _, err := NewGenerator(time.Now().Add(-time.Hour), uint16(1024)); !errors.Is(err, ErrInvalidMachineID) {
		t.Fatalf("using machine ID greater than 1023 should've returned error %v but got %v", ErrInvalidMachineID, err)
	}
}

func TestNextID(t *testing.T) {
	epoch := time.Now().Add(-time.Hour)
	machineID := uint16(123)

	g, err := NewGenerator(epoch, machineID)
	if err != nil {
		t.Fatalf("failed to create new generator: %v", err)
	}

	var j uint16
	var prevTs time.Time = time.Now().Add(-time.Hour)
	for range 1000 {
		id := g.NextID()

		if gotMachineID := g.MachineID(id); gotMachineID != machineID {
			t.Errorf("wrong machine ID. got %d, want %d", gotMachineID, machineID)
		}

		gotTs := g.Time(id)

		if gotTs.Before(prevTs) {
			t.Fatal("time should be greater than or equal to the previous time")
		}

		prevTs = gotTs

		gotSequence := g.Sequence(id)

		if gotSequence == 0 {
			j = 0
		}

		if j != gotSequence {
			t.Fatalf("wrong sequence. got %d, want %d", gotSequence, j)
		}

		j++
	}

}

func BenchmarkNextID(b *testing.B) {
	epoch := time.Now()
	machineID := uint16(0)

	g, err := NewGenerator(epoch, machineID)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		g.NextID()
	}
}
