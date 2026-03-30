package benchmark

import (
	"testing"
	"time"

	"github.com/16point7/malenkyi"
	"github.com/bwmarrin/snowflake"
	gosnowflake "github.com/godruoyi/go-snowflake"
	"github.com/sony/sonyflake"
)

func BenchmarkMalenkyi(b *testing.B) {
	g, err := malenkyi.NewGenerator(time.Now().Add(-time.Hour), 0)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		g.NextID()
	}
}

func BenchmarkMalenkyiParallel(b *testing.B) {
	g, err := malenkyi.NewGenerator(time.Now().Add(-time.Hour), 0)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			g.NextID()
		}
	})
}

func BenchmarkBwmarrin(b *testing.B) {
	node, err := snowflake.NewNode(0)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		node.Generate()
	}
}

func BenchmarkBwmarrinParallel(b *testing.B) {
	node, err := snowflake.NewNode(0)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			node.Generate()
		}
	})
}

func BenchmarkGodruoyi(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		gosnowflake.ID()
	}
}

func BenchmarkGodruoyiParallel(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gosnowflake.ID()
		}
	})
}

func BenchmarkSonyflake(b *testing.B) {
	sf, err := sonyflake.New(sonyflake.Settings{})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		sf.NextID()
	}
}

func BenchmarkSonyflakeParallel(b *testing.B) {
	sf, err := sonyflake.New(sonyflake.Settings{})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sf.NextID()
		}
	})
}
