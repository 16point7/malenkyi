package malenkyi

import (
	"testing"
	"time"
)

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
