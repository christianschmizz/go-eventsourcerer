package main

import (
	"errors"
	"fmt"

	es "github.com/christianschmizz/go-eventsourcerer"
	"github.com/rs/zerolog/log"
)

var accNumber = 10000

type accountCommandHandler struct {
	*es.BaseCommandHandler
	repo es.Repository
}

func NewAccountCommandHandler(repo es.Repository) *accountCommandHandler {
	return &accountCommandHandler{
		repo: repo,
	}
}

func (h *accountCommandHandler) Handle(c es.Command) error {
	log.Debug().Msgf("handling: %#v", c)

	var err error
	switch cmd := c.(type) {
	case OpenAccountCommand:
		err = h.handleOpenAccountCommand(cmd)
	case CloseAccountCommand:
		err = h.handleCloseAccountCommand(cmd)
	case FreezeAccountCommand:
		err = h.handleFreezeAccountCommand(cmd)
	case UnfreezeAccountCommand:
		err = h.handleUnfreezeAccountCommand(cmd)
	case DepositMoneyCommand:
		err = h.handleDepositMoneyCommand(cmd)
	case WithdrawMoneyCommand:
		err = h.handleWithdrawMoneyCommand(cmd)
	default:
		log.Fatal().Msgf("unknown command: %#v", cmd)
	}

	return err
}

func (h *accountCommandHandler) getAccount(n AccountNumber) (*Account, error) {
	accountMaybe, err := h.repo.GetByID(es.AggregateID(n))
	if err != nil {
		return nil, fmt.Errorf("account not found: %d", n)
	}
	account := accountMaybe.(*Account)
	return account, nil
}

func (h *accountCommandHandler) handleOpenAccountCommand(cmd OpenAccountCommand) error {
	if len(cmd.Username) == 0 {
		return cmd.Reject(fmt.Errorf("empty username is not allowed"))
	}

	accNumber += 1

	account := OpenAccount(cmd.Username, accNumber)
	h.repo.Save(account, es.InitialVersion)
	return nil
}

func (h *accountCommandHandler) handleCloseAccountCommand(cmd CloseAccountCommand) error {
	account, err := h.getAccount(cmd.AccountNumber)
	if err != nil {
		return err
	}

	if account.IsFrozen() {
		return cmd.Reject(fmt.Errorf("no withdrawl possible as account was frozen on %s", account.FrozenAt))
	}

	if account.Balance > 0 {
		return cmd.Reject(fmt.Errorf("account cannot be closed as there is still balance: %d", account.Balance))
	}

	account.Close()

	h.repo.Save(account, cmd.OriginalVersion)
	return nil
}

func (h *accountCommandHandler) handleDepositMoneyCommand(cmd DepositMoneyCommand) error {
	account, err := h.getAccount(cmd.AccountNumber)
	if err != nil {
		return err
	}

	if err := account.Deposit(cmd.Amount); err != nil {
		return cmd.Reject(fmt.Errorf("account not found"))
	}

	h.repo.Save(account, cmd.OriginalVersion)
	return nil
}

func (h *accountCommandHandler) handleWithdrawMoneyCommand(cmd WithdrawMoneyCommand) error {
	account, err := h.getAccount(cmd.AccountNumber)
	if err != nil {
		return err
	}

	if account.IsClosed() {
		return cmd.Reject(fmt.Errorf("no withdrawl possible as account was closed on %s", account.ClosedAt))
	}

	if account.IsFrozen() {
		return cmd.Reject(fmt.Errorf("no withdrawl possible as account was frozen on %s", account.FrozenAt))
	}

	resultingBalance := account.Balance - cmd.Amount
	if -resultingBalance > account.Limit {
		return cmd.Reject(fmt.Errorf("withdrawl would result in %d balance exceeding limit: %d", resultingBalance, account.Limit))
	}

	if err := account.Withdraw(cmd.Amount); err != nil {
		return cmd.Reject(errors.New("account not found"))
	}

	h.repo.Save(account, cmd.OriginalVersion)
	return nil
}

func (h *accountCommandHandler) handleFreezeAccountCommand(cmd FreezeAccountCommand) error {
	account, err := h.getAccount(cmd.AccountNumber)
	if err != nil {
		return err
	}

	account.Freeze()

	h.repo.Save(account, cmd.OriginalVersion)
	return nil
}

func (h *accountCommandHandler) handleUnfreezeAccountCommand(cmd UnfreezeAccountCommand) error {
	account, err := h.getAccount(cmd.AccountNumber)
	if err != nil {
		return err
	}

	account.Unfreeze()

	h.repo.Save(account, cmd.OriginalVersion)
	return nil
}
