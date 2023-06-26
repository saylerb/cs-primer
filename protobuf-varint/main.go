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
	var mask uint64 = 0x7f // 0b01111111 least significant 7 bits
	buffer := new(bytes.Buffer)

	for input > 0 {
		lowestSevenBits := uint8(input & mask)
		if input > 0x7F {
			withContinuationBit := lowestSevenBits | 0x80
			binary.Write(buffer, binary.BigEndian, withContinuationBit)
		} else {
			binary.Write(buffer, binary.BigEndian, lowestSevenBits)
		}
		input = input >> 7
	}

	return buffer.Bytes()
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
	buffer := bytes.NewReader(theSlice)
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
