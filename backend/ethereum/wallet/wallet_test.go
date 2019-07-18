// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package wallet

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"perun.network/go-perun/wallet/test"
)

const (
	keyDir      = "testdata"
	password    = "secret"
	sampleAddr  = "0x1234560000000000000000000000000000000000"
	invalidAddr = "0x12345600000000000000000000000000000000001"
	dataToSign  = "SomeLongDataThatShouldBeSignedPlease"

	keystoreAddr = "0xf4c288068b32474dedc3620233c"
	keyStorePath = "UTC--2019-06-07T12-12-48.775026092Z--3c5a96ff258b1f4c288068b32474dedc3620233c"
)

func TestGenericWalletTests(t *testing.T) {
	t.Parallel()
	setup := newTestSetup()
	test.GenericWalletTest(t, setup)
}

func TestGenericSignatureTests(t *testing.T) {
	t.Parallel()
	setup := newTestSetup()
	test.GenericSignatureTest(t, setup)
}

func TestGenericAddressTests(t *testing.T) {
	t.Parallel()
	setup := newTestSetup()
	test.GenericAddressTest(t, setup)
}

func TestAddress(t *testing.T) {
	t.Parallel()
	w := connectTmpKeystore(t)
	perunAcc := w.Accounts()[0]
	ethAcc := new(accounts.Account)

	unsetAccount := new(Account)
	nilAddr := common.BytesToAddress(make([]byte, 40, 40))

	assert.Equal(t, nilAddr.Bytes(), unsetAccount.Address().Bytes(), "Unset address should be nil")
	ethAcc.Address.SetBytes(perunAcc.Address().Bytes())
	assert.Equal(t, ethAcc.Address.Bytes(), perunAcc.Address().Bytes(), "Bytes should return same value as internal structure")
	assert.NotEqual(t, nil, ethAcc.Address.Bytes(), "Set address should not be nil")
}

func TestKeyStore(t *testing.T) {
	t.Parallel()
	w := new(Wallet)
	assert.NotNil(t, w.Connect("", ""), "Expected connect to fail")

	w = connectTmpKeystore(t)

	unsetAccount := new(Account)
	assert.False(t, w.Contains(unsetAccount), "Keystore should not contain empty account")
}

func TestBackend(t *testing.T) {
	t.Parallel()
	backend := new(Backend)
	addr, err := backend.NewAddressFromString(sampleAddr)

	assert.Nil(t, err, "Conversion of valid address should work")
	_, err = backend.NewAddressFromBytes(addr.Bytes())
	assert.Nil(t, err, "Conversion of valid address should work")
	_, err = backend.NewAddressFromBytes([]byte(invalidAddr))
	assert.NotNil(t, err, "Conversion from wrong address should fail")
	_, err = backend.NewAddressFromString(invalidAddr)
	assert.NotNil(t, err, "Conversion from wrong address should fail")
}

func newTestSetup() *test.Setup {
	return &test.Setup{
		Wallet:     new(Wallet),
		Path:       "./" + keyDir,
		WalletPW:   password,
		AccountPW:  password,
		Backend:    new(Backend),
		AddrString: sampleAddr,
		DataToSign: []byte(dataToSign),
	}
}

func connectTmpKeystore(t *testing.T) *Wallet {
	w := new(Wallet)
	assert.Nil(t, w.Connect(keyDir, password), "Unable to open keystore")
	assert.NotEqual(t, len(w.Accounts()), 0, "Wallet contains no accounts")
	return w
}

// Benchmarking part starts here

func BenchmarkGenericAccount(b *testing.B) {
	setup := newBenchSetup()
	test.GenericAccountBenchmark(b, setup)
}

func BenchmarkGenericWallet(b *testing.B) {
	setup := newBenchSetup()
	test.GenericWalletBenchmark(b, setup)
}

func BenchmarkGenericBackend(b *testing.B) {
	setup := newBenchSetup()
	test.GenericBackendBenchmark(b, setup)
}

func newBenchSetup() *test.Setup {
	// Filled with the same data as the testing
	return &test.Setup{
		Wallet:     new(Wallet),
		Path:       "./" + keyDir,
		WalletPW:   password,
		AccountPW:  password,
		Backend:    new(Backend),
		AddrString: sampleAddr,
		DataToSign: []byte(dataToSign),
	}
}
