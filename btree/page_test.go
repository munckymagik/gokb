package btree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPageFind(t *testing.T) {
	t.Run("when the page is empty", func(t *testing.T) {
		// Given
		p := newPage[int, emptyValue](nil)

		//  When
		entry, leaf := p.find(1)

		// Then
		require.Nil(t, entry)
		require.Equal(t, p, leaf)
	})
	t.Run("when the page has entries", func(t *testing.T) {
		// Given
		p := newPage[int, emptyValue](nil)

		//  When
		entry, leaf := p.find(1)

		// Then
		require.Nil(t, entry)
		require.Equal(t, p, leaf)
	})
	t.Run("when the page has children", func(t *testing.T) {})
}

type emptyValue = struct{}
