package main

import (
	"bytes"
	"fmt"
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

	t.Run("be able to convert a byte array less than 8 bytes to a 64 bit integer", func(t *testing.T) {
		got := BtoI([]byte{0x96})
		var want uint64 = 150
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
	t.Run("can encode 300 to protobuf varint", func(t *testing.T) {
		got := Encode(300)
		want := []byte{0xac, 0x02}
		compareByteSlices(t, got, want)
	})
	t.Run("can encode 128 to protobuf varint", func(t *testing.T) {
		got := Encode(128)
		want := []byte{0x80, 0x01}
		compareByteSlices(t, got, want)
	})
}

func TestDecodingVarint(t *testing.T) {
	t.Run("can decode the number 150", func(t *testing.T) {
		got := Decode([]byte{0x96, 0x01})
		var want uint64 = 150
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("can decode the number 1", func(t *testing.T) {
		got := Decode([]byte{0x01})
		var want uint64 = 1
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("can decode max int 1", func(t *testing.T) {
		got := Decode([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01})
		var want uint64 = 18446744073709551615
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("can decode the number 300", func(t *testing.T) {
		got := Decode([]byte{0xac, 0x02})
		var want uint64 = 300
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("can decode the number 128", func(t *testing.T) {
		got := Decode([]byte{0x80, 0x01})
		var want uint64 = 128
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

}

func TestRoundTrip(t *testing.T) {
	t.Run("can roundtrip encode/decode numbers 1 to ~1 billion", func(t *testing.T) {
		// currently this takes a while ~ 2 minutes
		var end uint64 = 1 << 30
		var i uint64
		fmt.Printf("roundtripping with numbers 1 to %v", end)
		for i = 1; i <= end; i++ {
			encoded := Encode(i)
			decoded := Decode(encoded)

			if i != decoded {
				fmt.Printf("error with %v", i)

				t.Errorf("got %v want %v", decoded, i)
			}
		}
	})
}

func compareByteSlices(t testing.TB, got, want []byte) {
	t.Helper()
	if !bytes.Equal(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}
