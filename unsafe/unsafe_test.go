package unsafe

import (
	"testing"
	stdunsafe "unsafe"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("fruits can be created as handles to the base type", func(t *testing.T) {
		// Given
		a := create(kOrange)
		b := create(kBanana)

		// Then
		require.Equal(t, kOrange, a.kind)
		require.Equal(t, kBanana, b.kind)
	})

	t.Run("a pointer of the base type uses less memory than interface{}", func(t *testing.T) {
		// Given
		var asBaseP *fruit
		var asIfP interface{} = &orange{}

		// Then
		require.Equal(t, uintptr(8), stdunsafe.Sizeof(asBaseP))
		require.Equal(t, uintptr(16), stdunsafe.Sizeof(asIfP))
	})

	t.Run("base type pointers can be downcasted", func(t *testing.T) {
		// Given
		a := create(kOrange)
		b := create(kBanana)

		// Then
		require.Nil(t, a.AsBanana())
		require.Nil(t, b.AsOrange())

		// When
		aa := a.AsOrange()
		bb := b.AsBanana()

		// Then
		require.Equal(t, kOrange, aa.kind)
		require.Equal(t, "orange", aa.name)

		require.Equal(t, kBanana, bb.kind)
		require.Equal(t, cYellow, bb.colour)
	})

	t.Run("method calls are not polymorphic", func(t *testing.T) {
		// Given
		a := create(kOrange)
		b := create(kBanana)
		aa := a.AsOrange()
		bb := b.AsBanana()

		// Then
		require.Equal(t, "I am fruit orange", a.Say())
		require.Equal(t, "I am fruit banana", b.Say())
		require.Equal(t, "I am orange", aa.Say())
		require.Equal(t, "I am banana", bb.Say())
	})
}
