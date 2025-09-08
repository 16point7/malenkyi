# malenkyi

маленький (malénʹkyj) adj. little, small, or tiny, often carrying an affectionate or tender nuance.

## Description
Lock-free, 64-bit ID generator. Supports up to 1024 machines generating 4096 IDs every milisecond for ~69.7 years.

## How to use
```go
g, _ := malenkyi.NewGenerator(epoch, machineID)

id := g.NextID() // 4506041370070220800
g.Time(id)       // 2025-09-07 21:19:40.825 -0700 PDT
g.MachineID(id)  // 0
g.Sequence(id)   // 0

id = g.NextID()  // 4506041370070220801
g.Time(id)       // 2025-09-07 21:19:40.825 -0700 PDT
g.MachineID(id)  // 0
g.Sequence(id)   // 1
```

## Implementation details

IDs are `int64`.

Bit 63: Signed bit. Always zero.
Bits 62-22: Time since epoch in miliseconds. Capacity is ~69.7 years.
Bits 21-12: Machine ID. Range is [0,1023].
Bits 11-0: Sequence number. Range is [0,4095].

Calling `NextID()` will panic if time has overflowed the capacity (41 bits).

Thread safety is guaranteed by a single compare-and-swap operation on an `int64` representing the previous ID. Bits 52-12 hold the time since epoch in miliseconds and bits 11-0 hold the sequence number.