// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package sim

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"io"

	"github.com/pkg/errors"

	"perun.network/go-perun/log"
	perun "perun.network/go-perun/wallet"
)

var curve = elliptic.P256()

// Backend implements the utility interface defined in the wallet package.
type Backend struct{}

// NewAddressFromString creates a new address from a string.
// DEPRECATED
func (h *Backend) NewAddressFromString(s string) (perun.Address, error) {
	return h.NewAddressFromBytes([]byte(s))
}

// NewAddressFromBytes creates a new address from a byte array.
// DEPRECATED
func (h *Backend) NewAddressFromBytes(data []byte) (perun.Address, error) {
	return h.DecodeAddress(bytes.NewReader(data))
}

// DecodeAddress decodes an address from the given Reader
func (h *Backend) DecodeAddress(r io.Reader) (perun.Address, error) {
	var addr Address
	return &addr, addr.Decode(r)
}

// VerifySignature verifies if a signature was made by this account.
func (h *Backend) VerifySignature(msg, sig []byte, a perun.Address) (bool, error) {
	addr, ok := a.(*Address)
	if !ok {
		log.Panic("Wrong address type passed to Backend.VerifySignature")
	}
	pk := (*ecdsa.PublicKey)(addr)

	r, s, err := deserializeSignature(sig)
	if err != nil {
		return false, errors.WithMessage(err, "could not deserialize signature")
	}

	// escda.Verify needs a digest as input
	// ref https://golang.org/pkg/crypto/ecdsa/#Verify
	return ecdsa.Verify(pk, digest(msg), r, s), nil
}

func digest(msg []byte) []byte {
	digest := sha256.Sum256(msg)
	return digest[:]
}
