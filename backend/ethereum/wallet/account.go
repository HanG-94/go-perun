// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package wallet

import (
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/pkg/errors"
	"perun.network/go-perun/log"
	perun "perun.network/go-perun/wallet"
)

// Account represents an ethereum account.
type Account struct {
	address Address
	Account *accounts.Account
	wallet  *Wallet
	locked  bool
	mu      sync.RWMutex
}

// Address returns the ethereum address of this account.
func (a *Account) Address() perun.Address {
	return &a.address
}

// Unlock unlocks this account.
func (a *Account) Unlock(password string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	err := a.wallet.Ks.Unlock(*a.Account, password)
	if err != nil {
		return err
	}
	a.locked = false
	return nil
}

// IsLocked checks if this account is locked.
func (a *Account) IsLocked() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.locked
}

// Lock locks this account.
func (a *Account) Lock() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	err := a.wallet.Ks.Lock(a.address.Address)
	if err != nil {
		return err
	}
	a.locked = true
	return nil
}

// SignData is used to sign data with this account.
func (a *Account) SignData(data []byte) ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	hash := prefixedHash(data)
	sig, err := a.wallet.Ks.SignHash(*a.Account, hash)
	if err != nil {
		return nil, errors.WithMessage(err, "could not sign data")
	}
	if sig[64] <= 1 {
		sig[64] = sig[64] + 27
	}
	log.Debug(sig)
	return sig, nil
}

// SignDataWithPW is used to sign a hash with this account and a pw.
func (a *Account) SignDataWithPW(password string, data []byte) ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	hash := prefixedHash(data)
	sig, err := a.wallet.Ks.SignHashWithPassphrase(*a.Account, password, hash)
	if err != nil {
		return nil, errors.WithMessage(err, "could not sign data")
	}
	if sig[64] <= 1 {
		sig[64] = sig[64] + 27
	}
	return sig, nil
}

func NewAccountFromEth(wallet *Wallet, account *accounts.Account) *Account {
	return &Account{
		address: Address{account.Address},
		Account: account,
		wallet:  wallet,
		locked:  true,
	}
}
