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
