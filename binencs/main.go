package main

import (
	"encoding/ascii85"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	data := []byte("12345678123456781234567812345678")
	dataSize := len(data)

	dataHex := hex.EncodeToString(data)
	dataBase32 := base32.StdEncoding.EncodeToString(data)
	dataBase64 := base64.StdEncoding.EncodeToString(data)
	dataAscii85 := make([]byte, ascii85.MaxEncodedLen(len(data)))
	_ = ascii85.Encode(dataAscii85, data)

	report("Input", string(data), dataSize)
	report("ascii85", string(dataAscii85), dataSize)
	report("Base64", dataBase64, dataSize)
	report("Base32", dataBase32, dataSize)
	report("Hex", dataHex, dataSize)

	decoded, err := decodeASCII85(dataSize, string(dataAscii85))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Decoded from ascii85: %s\n", string(decoded))
}

func report(name string, encoded string, before int) {
	fmt.Printf(
		"%-8s (len: %2d, ratio %f) %v \n",
		name,
		len(encoded),
		ratio(before, len(encoded)),
		encoded,
	)
}

func ratio(before, after int) float64 {
	return float64(after) / float64(before)
}

func decodeASCII85(decodedSize int, ascii85Encoded string) ([]byte, error) {
	// ascii85.Decode needs the destination buffer size to be the next
	// smallest multiple of 4 greater than the expected decoded size
	destBuffer := make([]byte, nextMultipleOf4(decodedSize))
	ndst, nsrc, err := ascii85.Decode(destBuffer, []byte(ascii85Encoded), true)

	if err != nil {
		return nil, err
	}
	if nsrc != len(ascii85Encoded) {
		return nil, fmt.Errorf("did not consume entire length of encoded message")
	}
	if ndst != decodedSize {
		return nil, fmt.Errorf("message did not decode to the expected size")
	}

	return destBuffer[0:ndst], nil
}

// nextMultipleOf4 rounds up value to the next smallest multiple of 4
func nextMultipleOf4(value int) int {
	return ((value + 3) / 4) * 4
}
