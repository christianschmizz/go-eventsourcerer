package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenBankAccount(t *testing.T) {
	account := OpenAccount("John Doe", 12345)
	assert.Len(t, account.Journal().SliceAndDrain(), 1)
	assert.Equal(t, "John Doe", account.AccountHolder)
	assert.Equal(t, AccountNumber(12345), account.Number)
	assert.Equal(t, 0, account.Balance)
	assert.Equal(t, 0, account.Limit)
	assert.True(t, account.ClosedAt.IsZero())
}

func TestBankAccount_Deposit(t *testing.T) {
	account := OpenAccount("John Doe", 12345)

	t.Run("deposit money", func(t *testing.T) {
		err := account.Deposit(100)
		assert.NoError(t, err)
		assert.Len(t, account.Journal().SliceAndDrain(), 2)
		assert.Equal(t, 100, account.Balance)
	})
}

func TestBankAccount_Withdraw(t *testing.T) {
	account := OpenAccount("John Doe", 12345)

	t.Run("deposit money", func(t *testing.T) {
		err := account.Deposit(100)
		assert.NoError(t, err)
		assert.Len(t, account.Journal().SliceAndDrain(), 2)
		assert.Equal(t, 100, account.Balance)
	})

	t.Run("withdraw money", func(t *testing.T) {
		err := account.Withdraw(12)
		assert.NoError(t, err)
		assert.Len(t, account.Journal().SliceAndDrain(), 1)
		assert.Equal(t, 88, account.Balance)
	})
}

func TestBankAccount_Close(t *testing.T) {
	account := OpenAccount("John Doe", 12345)
	assert.Equal(t, "John Doe", account.AccountHolder)
	assert.Equal(t, AccountNumber(12345), account.Number)
	assert.Equal(t, 0, account.Balance)
	assert.Equal(t, 0, account.Limit)
	assert.True(t, account.ClosedAt.IsZero())

	err := account.Close()
	assert.NoError(t, err)
	assert.False(t, account.ClosedAt.IsZero())

	assert.Len(t, account.Journal().SliceAndDrain(), 2)
}