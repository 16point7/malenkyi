package malenkyi

import (
	"errors"
	"slices"
	"sync"
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

func FuzzNextID(f *testing.F) {
	f.Add(100, time.Now().Add(-time.Hour).UnixMilli(), uint16(144))
	f.Add(1000, time.Now().Add(-24*365*7*time.Hour).UnixMilli(), uint16(1023))

	const numCallsPerGoroutine = 1000

	f.Fuzz(func(t *testing.T, numGoroutines int, epochMs int64, machineID uint16) {
		if numGoroutines <= 0 || numGoroutines > 1000 {
			numGoroutines = 100
		}

		nowMs := time.Now().UnixMilli()
		if epochMs > nowMs {
			return
		}

		if machineID > maxMachineID {
			return
		}

		g, err := NewGenerator(time.UnixMilli(epochMs), machineID)
		if err != nil {
			t.Fatalf("failed to create generator: %v\n", err)
		}

		wg := sync.WaitGroup{}
		wg.Add(numGoroutines)

		resultChan := make(chan []int64, numGoroutines)

		for range numGoroutines {
			go func() {
				results := make([]int64, 0, numCallsPerGoroutine)
				for range numCallsPerGoroutine {
					results = append(results, g.NextID())
				}
				resultChan <- results
				wg.Done()
			}()
		}

		wg.Wait()
		close(resultChan)

		results := make([]int64, 0, numGoroutines*numCallsPerGoroutine)
		for result := range resultChan {
			results = append(results, result...)
		}

		slices.Sort(results)

		var prevSeq uint16
		var prevTs time.Time
		for _, id := range results {
			gotTs, gotMachineID, gotSeq := g.Decompose(id)

			if gotMachineID != machineID {
				t.Fatalf("wrong machine id. got %d, want %d", gotMachineID, machineID)
			}

			if gotSeq == 0 {
				if !gotTs.After(prevTs) {
					t.Fatal("sequence reset should only happen when time is advanced")
				}
			} else {
				if !gotTs.Equal(prevTs) {
					t.Fatal("sequence advance should only happen in same time window")
				}
				if gotSeq != prevSeq+1 {
					t.Fatalf("sequence should only advance by 1. got %d, want %d", gotSeq, prevSeq+1)
				}
			}
			prevSeq = gotSeq
			prevTs = gotTs
		}
	})
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
