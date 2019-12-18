// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package channel

import (
	"context"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	wallettest "perun.network/go-perun/wallet/test"

	"perun.network/go-perun/backend/ethereum/channel/test"
	channeltest "perun.network/go-perun/channel/test"
	perunwallet "perun.network/go-perun/wallet"

	ethwallettest "perun.network/go-perun/backend/ethereum/wallet/test"
)

const nodeURL = "ws://localhost:8545"

const (
	keyDir   = "../wallet/testdata"
	password = "secret"

	keystoreAddr = "0x3c5A96FF258B1F4C288068B32474dedC3620233c"
	keyStorePath = "UTC--2019-06-07T12-12-48.775026092Z--3c5a96ff258b1f4c288068b32474dedc3620233c"
)

const timeout = 300 * time.Millisecond

func TestFunder_Fund(t *testing.T) {
	f := newSimulatedFunder()
	assert.Panics(t, func() { f.Fund(context.Background(), channel.FundingReq{}) }, "Funding with invalid funding req should fail")
	req := channel.FundingReq{
		Params:     &channel.Params{},
		Allocation: &channel.Allocation{},
		Idx:        0,
	}
	assert.NoError(t, f.Fund(context.Background(), req), "Funding with no assets should succeed")
	parts := []perunwallet.Address{
		&wallet.Address{Address: f.account.Address},
	}
	rng := rand.New(rand.NewSource(1337))
	app := channeltest.NewRandomApp(rng)
	params := channel.NewParamsUnsafe(uint64(0), parts, app.Def(), big.NewInt(rng.Int63()))
	allocation := newValidAllocation(parts, f.ethAssetHolder)
	req = channel.FundingReq{
		Params:     params,
		Allocation: allocation,
		Idx:        0,
	}
	// Test with valid context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	assert.NoError(t, f.Fund(ctx, req), "funding with valid request should succeed")
	assert.NoError(t, f.Fund(ctx, req), "multiple funding should succeed")
	cancel()
	// Test already closed context
	assert.Error(t, f.Fund(ctx, req), "funding with already cancelled context should fail")
}

func TestFunder_Fund_multi(t *testing.T) {
	t.Run("1-party funding", func(t *testing.T) { testFunderFunding(t, 1) })
	t.Run("2-party funding", func(t *testing.T) { testFunderFunding(t, 2) })
}

func testFunderFunding(t *testing.T, n int) {
	rng := rand.New(rand.NewSource(1337))
	ctx, cancel := context.WithTimeout(context.Background(), 10*timeout)
	defer cancel()

	simBackend := test.NewSimulatedBackend()
	ks := ethwallettest.GetKeystore()
	deployAccount := wallettest.NewRandomAccount(rng).(*wallet.Account).Account
	simBackend.FundAddress(ctx, deployAccount.Address)
	contractBackend := NewContractBackend(simBackend, ks, deployAccount)
	// Deploy Assetholder
	assetETH, err := DeployETHAssetholder(ctx, contractBackend, deployAccount.Address)
	require.NoError(t, err)
	t.Log(assetETH.String())

	parts := make([]perunwallet.Address, n)
	funders := make([]*Funder, n)
	for i := 0; i < n; i++ {
		acc := wallettest.NewRandomAccount(rng).(*wallet.Account)
		simBackend.FundAddress(ctx, acc.Account.Address)
		parts[i] = acc.Address()
		cb := NewContractBackend(simBackend, ks, acc.Account)
		funders[i] = NewSimulatedFunder(cb, assetETH)
	}
	app := channeltest.NewRandomApp(rng)
	params := channel.NewParamsUnsafe(uint64(0), parts, app.Def(), big.NewInt(rng.Int63()))
	allocation := newValidAllocation(parts, assetETH)
	for i, funder := range funders {
		req := channel.FundingReq{
			Params:     params,
			Allocation: allocation,
			Idx:        uint16(i),
		}
		err := funder.Fund(ctx, req)
		assert.NoError(t, err, "funding should succeed")
	}
}

func newSimulatedFunder() *Funder {
	// Set KeyStore
	wall := new(wallet.Wallet)
	wall.Connect(keyDir, password)
	acc := wall.Accounts()[0].(*wallet.Account)
	acc.Unlock(password)
	ks := wall.Ks
	simBackend := test.NewSimulatedBackend()
	simBackend.FundAddress(context.Background(), acc.Account.Address)
	cb := ContractBackend{simBackend, ks, acc.Account}
	// Deploy Assetholder
	assetETH, err := DeployETHAssetholder(context.Background(), cb, acc.Account.Address)
	if err != nil {
		panic(err)
	}
	return &Funder{
		ContractBackend: cb,
		ethAssetHolder:  assetETH,
	}
}

func newValidAllocation(parts []perunwallet.Address, assetETH common.Address) *channel.Allocation {
	// Create assets slice
	assets := []channel.Asset{
		&Asset{Address: assetETH},
	}
	rng := rand.New(rand.NewSource(1337))
	ofparts := make([][]channel.Bal, len(parts))
	for i := 0; i < len(ofparts); i++ {
		ofparts[i] = make([]channel.Bal, len(assets))
		for k := 0; k < len(assets); k++ {
			// create new random balance in range [1,999]
			ofparts[i][k] = big.NewInt(rng.Int63n(999) + 1)
		}
	}
	return &channel.Allocation{
		Assets:  assets,
		OfParts: ofparts,
	}
}
