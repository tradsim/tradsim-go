package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOrderBook(t *testing.T) {

	require := require.New(t)

	book := NewOrderBook()

	require.NotNil(book)
	require.Empty(book.Symbols)
}
