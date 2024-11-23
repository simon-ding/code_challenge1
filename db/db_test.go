package db

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalanceToInt(t *testing.T) {
	b := big.NewRat(100, 1)
	i := balanceToInt(b)
	assert.Equal(t, int64(10000), i)
}
