package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func Hello() string {
	return "Hello, World!"
}

func main() {
	fmt.Println(Hello())
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
	var lastByte byte
	var moreBytes = true
	var currentByteIndex = 0
	for moreBytes {
		currentByte := bytes[currentByteIndex]
		//fmt.Printf("index: %v, bytes: %b\n", currentByteIndex, currentByte)
		mostSignificantBit := currentByte & 0x80
		//fmt.Printf("mask bit: %b\n", 0x80)
		//fmt.Printf("mostSignificant bit: %b\n", mostSignificantBit)
		if mostSignificantBit != 0x80 {
			moreBytes = false
			//fmt.Printf("setting moreBytes to %v\n", moreBytes)

		}
		//fmt.Printf("index: %v, bytes: %b\n", currentByteIndex, currentByte)
		withoutContinuation := currentByte & 0x7f // drop the continuation bit
		//fmt.Printf("after dropping continuation bit: %b\n", withoutContinuation)
		leastSignificantBit := withoutContinuation & 0x01
		//fmt.Printf("leastSignificantBit: %b\n", leastSignificantBit)
		leastSignificantBit = leastSignificantBit << 7
		//fmt.Printf("leastSignificantBit after shifting: %b\n", leastSignificantBit)
		if lastByte > 0 {
			combined := lastByte | leastSignificantBit
			//fmt.Printf("last byte: %b\n", lastByte)
			//fmt.Printf("combined after shifting: %b\n", combined)
			reversed = append(reversed, combined)
		} else if !moreBytes {
			reversed = append(reversed, withoutContinuation)
		}
		lastByte = withoutContinuation
		currentByteIndex = currentByteIndex + 1
	}

	//fmt.Printf("result: %b\n", reversed)

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
