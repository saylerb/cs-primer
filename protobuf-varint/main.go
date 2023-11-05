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
	buffer := new(bytes.Buffer)

	for input > 0 {
		lowestSevenBits := uint8(input & bitMask)
		if input > bitMask {
			var eighthBitOnMask uint8 = 0x80
			withContinuationBit := lowestSevenBits | eighthBitOnMask
			binary.Write(buffer, binary.BigEndian, withContinuationBit)
		} else {
			binary.Write(buffer, binary.BigEndian, lowestSevenBits)
		}
		input = input >> 7
	}
	return buffer.Bytes()
}

func Decode(bytes []byte) uint64 {
	var reversed []byte
	var previousByte byte
	var moreBytes = true
	var currentByteIndex = 0
	var continuationBitOnMask byte = 0x80
	for moreBytes {
		currentByte := bytes[currentByteIndex]
		mostSignificantBit := currentByte & continuationBitOnMask
		if mostSignificantBit != continuationBitOnMask {
			moreBytes = false
		}
		currentByteWithoutContinuationBit := currentByte & 0x7f                // set the continuation bit to zero
		currentLeastSignificantBit := currentByteWithoutContinuationBit & 0x01 // take the least significant bit
		if previousByte > 0 {
			// set the most significant bit on previous byte to value of current least significant bit
			combined := previousByte | (currentLeastSignificantBit << 7)
			reversed = append(reversed, combined)
		} else if !moreBytes {
			reversed = append(reversed, currentByteWithoutContinuationBit)
		}
		previousByte = currentByteWithoutContinuationBit
		currentByteIndex = currentByteIndex + 1
	}
	return BtoI(reversed)
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
