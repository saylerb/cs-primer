package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func main() {
	files := []string{"1.uint64", "150.uint64", "maxint.uint64"}
	for _, file := range files {
		bytes := ScanIntoByteSlice(file)
		int := BtoI(bytes)
		fmt.Printf("%v input as bytes: %x, input as decimal number: %v\n", file, bytes, int)
		encoded := Encode(int)
		fmt.Printf("%v encoded as protobuf varint: %x\n", file, encoded)
		decoded := Decode(encoded)
		fmt.Printf("%v varint decoded as uint64: %v\n", file, decoded)
		fmt.Println("------\n", file, decoded)
	}
}

func Encode(input uint64) []byte {
	var bitMask uint64 = 0x7f // 0b01111111 least significant 7 bits
	var buffer []byte

	for input > 0 {
		sevenLeastSignificantBits := uint8(input & bitMask)
		if input > bitMask {
			var eighthBitOnMask uint8 = 0x80
			withContinuationBit := sevenLeastSignificantBits | eighthBitOnMask
			buffer = append(buffer, withContinuationBit)
		} else {
			buffer = append(buffer, sevenLeastSignificantBits)
		}
		input = input >> 7
	}
	return buffer
}

func Decode(allBytes []byte) uint64 {
	var accumulator = 0
	var sevenLeastSignificantBitsMask uint8 = 0x7f // 0b01111111 least significant 7 bits

	for i := len(allBytes) - 1; i >= 0; i-- {
		accumulator <<= 7
		currentByte := allBytes[i]
		sevenLeastSignificantBits := currentByte & sevenLeastSignificantBitsMask
		accumulator += int(sevenLeastSignificantBits)
	}
	return uint64(accumulator)
}

func ReadBinaryFileToInteger(filename string) uint64 {
	return BtoI(ScanIntoByteSlice(filename))
}

func ItoB(num uint64) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, num)
	return buffer.Bytes()
}

func BtoI(theSlice []byte) uint64 {
	var result uint64

	var padded []byte
	for i := 0; i <= 7-len(theSlice); i++ {
		padded = append(padded, 0)
	}
	padded = append(padded, theSlice...)

	buffer := bytes.NewReader(padded)
	binary.Read(buffer, binary.BigEndian, &result)
	return result
}

func ScanIntoByteSlice(fileName string) []byte {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)

	var data []byte

	for scanner.Scan() {
		someBytes := scanner.Bytes()
		data = append(data, someBytes...)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
