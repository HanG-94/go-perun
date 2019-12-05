// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package channel

import (
	"encoding/hex"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/channel/test"
	perunwallet "perun.network/go-perun/wallet"
	wallettest "perun.network/go-perun/wallet/test"
)

func init() {
	channel.SetAppBackend(new(test.NoAppBackend))
}

func TestGenericTests(t *testing.T) {
	setup := newChannelSetup()
	test.GenericBackendTest(t, setup)
}

func newChannelSetup() *test.Setup {
	rng := rand.New(rand.NewSource(1337))

	app := test.NewRandomApp(rng)
	app2 := test.NewRandomApp(rng)

	params := test.NewRandomParams(rng, app.Def())
	params2 := test.NewRandomParams(rng, app2.Def())

	state := test.NewRandomState(rng, params)
	state2 := test.NewRandomState(rng, params2)
	state2.IsFinal = !state.IsFinal

	createAddr := func() perunwallet.Address {
		addr := wallet.NewRandomAddress(rng)
		return &addr
	}

	return &test.Setup{
		Params:        params,
		Params2:       params2,
		State:         state,
		State2:        state2,
		Account:       wallettest.NewRandomAccount(rng),
		RandomAddress: createAddr,
	}
}

func TestChannelID(t *testing.T) {
	tests := []struct {
		name        string
		aliceAddr   string
		bobAddr     string
		appAddr     string
		challengDur uint64
		nonceStr    string
		channelID   string
	}{
		{"Test case 1",
			"0xf17f52151EbEF6C7334FAD080c5704D77216b732",
			"0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef",
			"0x9FBDa871d559710256a2502A2517b794B482Db40",
			uint64(60),
			"B0B0FACE",
			"f27b90711d11d10a155fc8ba0eed1ffbf449cf3730d88c0cb77b98f61750ab34"},
		{"Test case 2",
			"0x0000000000000000000000000000000000000000",
			"0x0000000000000000000000000000000000000000",
			"0x0000000000000000000000000000000000000000",
			uint64(0),
			"0",
			"c8ac0e8f7eeea864a050a8626dfa0ffb916f43c90bc6b2ba68df6ed063c952e2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce, ok := new(big.Int).SetString(tt.nonceStr, 16)
			assert.True(t, ok, "Setting the nonce should not fail")
			alice, err := perunwallet.NewAddressFromString(tt.aliceAddr)
			assert.NoError(t, err, "Creating Alices address should not fail")
			bob, err := perunwallet.NewAddressFromString(tt.bobAddr)
			assert.NoError(t, err, "Creating Bobs address should not fail")
			app, err := perunwallet.NewAddressFromString(tt.appAddr)
			assert.NoError(t, err, "Creating the MockApp address should not fail")
			params := channel.Params{
				ChallengeDuration: tt.challengDur,
				Nonce:             nonce,
				Parts:             []perunwallet.Address{alice, bob},
				App:               test.NewMockApp(app),
			}
			cID := channel.ChannelID(&params)
			preCalc, err := hex.DecodeString(tt.channelID)
			assert.NoError(t, err, "Decoding the channelID should not error")
			assert.Equal(t, preCalc, cID[:], "ChannelID should match the testcase")
		})
	}
}
