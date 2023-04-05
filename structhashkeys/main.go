package main

import "fmt"

type travelTimeKey struct {
	zoneDrnID, restaurantDrnID, riderDrnID string

	// if we uncomment this it's no longer comparable literal to literal
	// ./prog.go:15:26: invalid map key type travelTimeKey
	// ./prog.go:34:35: invalid operation: key == key2 (struct containing []int cannot be compared)
	// makesNotComparable []int
}

func main() {
	fmt.Println("Hello, 世界")

	travelTimes := make(map[travelTimeKey]int)

	key := travelTimeKey{
		zoneDrnID:       "zone",
		restaurantDrnID: "restaurant",
		riderDrnID:      "rider",
	}

	travelTimes[key] = 13

	key2 := travelTimeKey{
		zoneDrnID:       "zone",
		restaurantDrnID: "restaurant",
		riderDrnID:      "rider",
	}

	fmt.Printf("&key = %p\n", &key)
	fmt.Printf("&key2 = %p\n", &key2)
	fmt.Printf("&key == &key2 = %v\n", &key == &key2)
	fmt.Printf("key == key2 = %v\n", key == key2)

	value, ok := travelTimes[key2]

	if !ok {
		fmt.Println("Missing")
	} else {
		fmt.Printf("Value = %d\n", value)
	}

}
