# Protobuf varint

Base 128 Varint
[varint](https://protobuf.dev/programming-guides/encoding/#varints) Protocol
Buffers.

### Context

- we know that the varint implementation in protobuf uses between 2 and 10 bytes to represent unsigned 64 bit
  integers
- this encoding uses 7 bits for the payload, and a single bit as continutation bit (indicating if there are more bytes
  following that represent that number)

### Encoding Algorithm

1. consider the number in chunks of 7 bits
    - save off the least significant 7 bits using bitmask
    - reserve the 8th bit (most significant) as continuation bit
2. add continuation bit as 8th bit if needed - if the number being processed is bigger than 127 (7 bits) there will be
   addition bytes to represent it, and the loop will continue
3. write this byte (with or without continuation bit) to a list of bytes that hold the result
4. shift the bits of the number to encode left by 7 bits
5. continue from beginning, until the number to encode reaches zero