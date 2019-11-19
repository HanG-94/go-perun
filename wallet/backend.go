// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package wallet

import "io"

// Sig is a single signature
type Sig = []byte

// backend is set to the global wallet backend. It must be set through
// backend.Set(Collection).
var backend Backend

// Backend provides useful methods for this blockchain.
type Backend interface {
	// NewAddressFromString creates a new address from the natural string representation of this blockchain.
	NewAddressFromString(s string) (Address, error)

	// NewAddressFromBytes creates a new address from a byte array.
	NewAddressFromBytes(data []byte) (Address, error)

	// DecodeAddress reads and decodes an address from an io.Writer
	DecodeAddress(io.Reader) (Address, error)

	// VerifySignature verifies if this signature was signed by this address.
	// It should return an error iff the signature or message are malformed.
	// If the signature does not match the address it should return false, nil
	VerifySignature(msg []byte, sign Sig, a Address) (bool, error)
}

// SetBackend sets the global wallet backend. Must not be called directly but
// through backend.Set().
func SetBackend(b Backend) {
	backend = b
}

// NewAddressFromString calls NewAddressFromString of the current backend
func NewAddressFromString(s string) (Address, error) {
	return backend.NewAddressFromString(s)
}

// NewAddressFromBytes calls NewAddressFromBytes of the current backend
func NewAddressFromBytes(data []byte) (Address, error) {
	return backend.NewAddressFromBytes(data)
}

// DecodeAddress calls DecodeAddress of the current backend
func DecodeAddress(r io.Reader) (Address, error) {
	return backend.DecodeAddress(r)
}

// VerifySignature calls VerifySignature of the current backend
func VerifySignature(msg []byte, sign Sig, a Address) (bool, error) {
	return backend.VerifySignature(msg, sign, a)
}
