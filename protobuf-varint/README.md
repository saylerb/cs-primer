# Protobuf varint

Implementation of the Base 128 [varint](https://protobuf.dev/programming-guides/encoding/#varints) encoding used in
Protocol Buffers.

### Context

- we know that the varint implementation in protobuf uses between 2 and 10 bytes to represent unsigned 64 bit
  integers
- if the number is less than a single byte (7 bits), we can represent the number with a single byte
- If the number is larger than that, we encode the numbers in a sequence of bytes, with each byte using 7 bits for the
  payload,
  single bit as continuation bit (indicating if there are more bytes
  following that represent that number)

### Input

There are three files (1.uint64, 150.uint64, and maxint.uint64), each file is a
single number represented as binary data encoded as an unsigned 64-bit big endian integer.

### Encoding

### Encoding Algorithm

1. consider the number in chunks of 7 bits
    - save off the least significant 7 bits using bitmask
    - reserve the 8th bit (most significant) as continuation bit
2. add continuation bit as 8th bit if needed - if the number being processed is bigger than 127 (7 bits) there will be
   addition bytes to represent it, and the loop will continue
3. write this byte (with or without continuation bit) to a list of bytes that hold the result
4. shift the bits of the number to encode left by 7 bits
5. continue from beginning, until the number to encode reaches zero

### Open Questions

Why do we convert to from big endian to little endian order, and vice versa,
when encoding and decoding, respectively?

### Observations

Using `binary.Write` seems to be much slower than using built in `append`, when encoding and collecting up the
binary data. Writing a [benchmark](https://pkg.go.dev/testing#hdr-Benchmarks), there is the comparison of the output:

`binary.Write`:

```bash
go test -bench=. -benchtime=20s -benchmem
goos: darwin
goarch: arm64
pkg: varint
BenchmarkRoundTrip1/can_roundtrip_encode/decode_numbers_1_to_1-10         	1000000000	         0.6244 ns/op	       0 B/op	       0 allocs/op
BenchmarkRoundTrip10000/can_roundtrip_encode/decode_numbers_1_to_10000-10 	   32196	    732676 ns/op	 1139766 B/op	   39869 allocs/op
BenchmarkRoundTrip10_000_000/can_roundtrip_encode/decode_numbers_1_to_10000000-10         	      22	1028888496 ns/op	1157893836 B/op	57886355 allocs/op
PASS
ok  	varint	55.982s
```

using `append`:

```bash
go test -bench=. -benchtime=20s -benchmem
goos: darwin
goarch: arm64
pkg: varint
BenchmarkRoundTrip1/can_roundtrip_encode/decode_numbers_1_to_1-10         	1000000000	         0.6256 ns/op	       0 B/op	       0 allocs/op
BenchmarkRoundTrip10000/can_roundtrip_encode/decode_numbers_1_to_10000-10 	  138283	    170116 ns/op	   79992 B/op	    9999 allocs/op
BenchmarkRoundTrip10_000_000/can_roundtrip_encode/decode_numbers_1_to_10000000-10         	     123	 193308290 ns/op	80000413 B/op	10000000 allocs/op
PASS
ok  	varint	69.745s
```

You can see for the 3rd Benchamrk Scenario for encoding numbers from 1 to 10 million:

- 1.1 billion bytes processed for `binary.Write` vs. 80 million for `append` (~1347% more bytes for `binary.Write`)
- there were 58 million memory allocations per operation for `binary.Write`, vs. 10 million when using `append` (~478%
  more allocations for `binary.Write`)

I believe this is due to `binary.Write` making an additional copy of the input data prior to writing to the buffer,
but more investigation is needed to understand why.
