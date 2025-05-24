// This file demonstrates using unsafe pointers to hold references to
// structs with a common base type, without having to use interface{}. The
// benefit being that the memory size of a pointer to a concrete type is half
// the size of interface{}.
// Technique seen in https://github.com/tidwall/rtree.

package unsafe

import "unsafe"

type kind int8

const (
	kNone kind = iota
	kOrange
	kBanana
)

func (k kind) String() string {
	switch k {
	case kOrange:
		return "orange"
	case kBanana:
		return "banana"
	default:
		return ""
	}
}

type colour uint32

const (
	cNone colour = iota
	cYellow
)

type fruit struct {
	kind kind
}
type orange struct {
	fruit
	name string
}
type banana struct {
	fruit
	colour colour
}

func (f *fruit) Say() string {
	return "I am fruit " + f.kind.String()
}
func (f *fruit) AsOrange() *orange {
	if f.kind != kOrange {
		return nil
	}
	return (*orange)(unsafe.Pointer(f))
}
func (f *fruit) AsBanana() *banana {
	if f.kind != kBanana {
		return nil
	}
	return (*banana)(unsafe.Pointer(f))
}
func (f *orange) Say() string {
	return "I am orange"
}
func (f *banana) Say() string {
	return "I am banana"
}

func create(kind kind) *fruit {
	switch kind {
	case kOrange:
		return (*fruit)(unsafe.Pointer(&orange{
			fruit: fruit{kind: kOrange},
			name:  "orange",
		}))
	case kBanana:
		return (*fruit)(unsafe.Pointer(&banana{
			fruit:  fruit{kind: kBanana},
			colour: cYellow,
		}))
	default:
		return nil
	}
}
