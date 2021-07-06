package main

import "fmt"

func main() {
	// https://blog.golang.org/strings
	const nihongo = "日本語"
	for index, runeValue := range nihongo {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
	}

	fmt.Println("~~~~~~~~~")
	const charset = "\x09\x0A\x0D\x20\uD7FF\uE000\uFFFD\U00010000\U0010FFFF"
	for _, v := range charset {
		identify(v)
	}

	// Note: \uD800, \uDFFF are invalid code points
	var invalid = []int32{0x08, 0x0B, 0x0C, 0x1F, 0xD800, 0xDFFF, 0xFFFE, 0xFFFF, 0x110000}
	for _, v := range invalid {
		identify(v)
	}
}

func identify(runeValue rune) {
	if inXMLCharset(runeValue) {
		fmt.Printf("%#U is in the charset\n", runeValue)
	} else {
		fmt.Printf("%#U is NOT in the charset\n", runeValue)
	}
}

func inXMLCharset(v rune) bool {
	// Valid characters in a SQS SendMessage body
	// 	#x9 | #xA | #xD | [#x20-#xD7FF] | [#xE000-#xFFFD] | [#x10000-#x10FFFF]
	// From https://www.w3.org/TR/REC-xml/#charsets
	switch {
	case v == '\x09', v == '\x0A', v == '\x0D':
		fallthrough
	case v >= '\x20' && v <= '\uD7FF':
		fallthrough
	case v >= '\uE000' && v <= '\uFFFD':
		fallthrough
	case v >= '\U00010000' && v <= '\U0010FFFF':
		return true
	default:
		return false
	}
}
