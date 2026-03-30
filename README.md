# malenkyi

маленький (malénʹkyj) *adj.* little, small, or tiny, often carrying an affectionate or tender nuance.

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

- Bit 63: Signed bit. Always zero.
- Bits 62-22: Time since epoch in milliseconds. Capacity is ~69.7 years.
- Bits 21-12: Machine ID. Range is [0,1023].
- Bits 11-0: Sequence number. Range is [0,4095].

Calling `NextID()` will panic if time has overflowed the capacity (41 bits).

Thread safety is guaranteed by a single compare-and-swap operation on an `int64` representing the previous ID. Bits 52-12 hold the time since epoch in milliseconds and bits 11-0 hold the sequence number.

## Benchmarks

Comparison against other Go snowflake libraries. All benchmarks run on an Intel i7-8700K (12 threads) with Go 1.24.

### Single-threaded

| Library | ns/op | allocs/op |
|---|---|---|
| malenkyi | 243.9 | 0 |
| [bwmarrin/snowflake](https://github.com/bwmarrin/snowflake) | 243.9 | 0 |
| [godruoyi/go-snowflake](https://github.com/godruoyi/go-snowflake) | 244.1 | 0 |
| [sony/sonyflake](https://github.com/sony/sonyflake) | 38848 | 0 |

### Parallel (GOMAXPROCS=12)

| Library | ns/op | allocs/op |
|---|---|---|
| malenkyi | 244.0 | 0 |
| [bwmarrin/snowflake](https://github.com/bwmarrin/snowflake) | 244.0 | 0 |
| [godruoyi/go-snowflake](https://github.com/godruoyi/go-snowflake) | 244.2 | 0 |
| [sony/sonyflake](https://github.com/sony/sonyflake) | 38836 | 0 |

malenkyi, bwmarrin, and godruoyi all share the same bit layout (41+10+12) and are bottlenecked by `time.Now()` (~244 ns) rather than their concurrency mechanism. sonyflake is ~160x slower because it uses 10ms time units with an 8-bit sequence (256 IDs per unit), sleeping on overflow.