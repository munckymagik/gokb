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
