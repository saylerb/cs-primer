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

func ReadBinaryFileToInteger(filename string) uint64 {
	return BtoI(ScanIntoByteSlice(filename))
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
