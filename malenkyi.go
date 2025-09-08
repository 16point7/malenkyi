package malenkyi

import (
	"errors"
	"sync/atomic"
	"time"
)

// ID schema:
//
//  41 bits: miliseconds since epoch, ~ 69.684 years                     12 bits: sequence, [0,4095]
//  *******   ********   ********   ********   ********   **             ****   ********
//  *                                                      *             *             *
// 00000000 - 00000000 - 00000000 - 00000000 - 00000000 - 00000000 - 00000000 - 00000000
//    *                                                        *           *
// 1 bit: signed bit, always zero                           ******   ****
//                                                             10 bits: machine ID, [0,1023]

// ID generator.
// Can be used concurrently without locks.
type Generator struct {
	machineID int64
	epochMs   int64
	prevseq   atomic.Int64
}

func NewGenerator(epoch time.Time, machineID uint16) (*Generator, error) {
	if machineID > maxMachineID {
		return nil, ErrInvalidMachineID
	}

	if epoch.After(time.Now()) {
		return nil, ErrInvalidEpoch
	}

	return &Generator{machineID: int64(machineID), epochMs: epoch.UnixMilli()}, nil
}

var (
	ErrInvalidMachineID = errors.New("machine ID must be in the range [0, 1023]")
	ErrInvalidEpoch     = errors.New("epoch cannot be in the future")
)

const (
	tsBits        = 41
	machineIDBits = 10
	seqBits       = 12

	maxTs        = (1 << tsBits) - 1
	maxMachineID = (1 << machineIDBits) - 1
	maxSequence  = (1 << seqBits) - 1
)

func (g *Generator) NextID() int64 {
	for {
		ts := time.Now().UnixMilli() - g.epochMs

		if ts > maxTs {
			panic("timestamp overflowed 41 bits")
		}

		seq, ok := g.nextSequence(ts)
		if !ok {
			continue
		}

		return (ts << (machineIDBits + seqBits)) | (g.machineID << seqBits) | seq
	}
}

func (g *Generator) nextSequence(nowMs int64) (int64, bool) {
	for {
		prevseq := g.prevseq.Load()

		prevMs := prevseq >> seqBits

		if nowMs > prevMs {
			if g.prevseq.CompareAndSwap(prevseq, nowMs<<seqBits) {
				return 0, true
			}
			continue
		}

		if nowMs == prevMs {
			prevSeq := prevseq & maxSequence

			if prevSeq >= maxSequence {
				return 0, false
			}

			if g.prevseq.CompareAndSwap(prevseq, (nowMs<<seqBits)|(prevSeq+1)) {
				return prevSeq + 1, true
			}

			continue
		}

		return 0, false
	}
}

func (g *Generator) Time(id int64) time.Time {
	return time.UnixMilli(g.epochMs).Add(time.Millisecond * time.Duration(id>>(machineIDBits+seqBits)))
}

func (g *Generator) MachineID(id int64) uint16 {
	return uint16((maxMachineID & id) >> seqBits)
}

func (g *Generator) Sequence(id int64) uint16 {
	return uint16(maxSequence & id)
}

func (g *Generator) Decompose(id int64) (ts time.Time, machineID, seq uint16) {
	return g.Time(id), g.MachineID(id), g.Sequence(id)
}
