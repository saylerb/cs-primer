package main

import (
	"bytes"
	"testing"
)

func TestReadBinaryFileToInteger(t *testing.T) {
	t.Run("test comparing two byte slices", func(t *testing.T) {
		one := []byte{1, 2, 3}
		two := []byte{1, 2, 3}

		compareByteSlices(t, one, two)
	})

	t.Run("read the binary number 1 from file", func(t *testing.T) {
		got := ScanIntoByteSlice("1.uint64")
		want := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1}
		compareByteSlices(t, got, want)
	})

	t.Run("be able to convert a byte array to integer", func(t *testing.T) {
		u64bitIntOne := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1}
		got := BtoI(u64bitIntOne)
		var want uint64 = 1

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("be able to convert integer to byte slice", func(t *testing.T) {
		var one uint64 = 1
		want := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1}
		got := ItoB(one)

		compareByteSlices(t, got, want)
	})

	t.Run("read uint64 one from binary file", func(t *testing.T) {
		got := ReadBinaryFileToInteger("1.uint64")
		var want uint64 = 1

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("read uint64 150 from binary file", func(t *testing.T) {
		got := ReadBinaryFileToInteger("150.uint64")
		var want uint64 = 150

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("read max unint64 from binary file", func(t *testing.T) {
		got := ReadBinaryFileToInteger("maxint.uint64")
		var want uint64 = 18446744073709551615
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestEncodingVarint(t *testing.T) {
	t.Run("can encode the number 1 to protobuf varint", func(t *testing.T) {
		got := Encode(ReadBinaryFileToInteger("1.uint64"))
		want := []byte{0x01}
		compareByteSlices(t, got, want)
	})

	t.Run("can encode the number 150 to protobuf varint", func(t *testing.T) {
		got := Encode(ReadBinaryFileToInteger("150.uint64"))
		want := []byte{0x96, 0x01}
		compareByteSlices(t, got, want)
	})
	t.Run("can encode maxint to protobuf varint", func(t *testing.T) {
		got := Encode(ReadBinaryFileToInteger("maxint.uint64"))
		want := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}
		compareByteSlices(t, got, want)
	})
}

func compareByteSlices(t testing.TB, got, want []byte) {
	t.Helper()
	if !bytes.Equal(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}
