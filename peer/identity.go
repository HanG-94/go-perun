// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package peer

import (
	"io"

	"github.com/pkg/errors"

	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire/msg"
)

// Identity is a node's permanent Perun identity, which is used to establish
// authenticity within the Perun peer-to-peer network. For now, it is just a
// stub.
type Identity = wallet.Account

// Authenticate runs an authentication protocol on a connection.
// The protocol exchanges Perun addresses and establishes authenticity. NOTE:
// this protocol is not secure against man-in-the-middle attacks.
//
// Authenticate() returns the peer's address, if successful, or an error.
func Authenticate(id Identity, conn Conn) (Address, error) {
	sent := make(chan error, 1)
	go func() { sent <- conn.Send(NewAuthResponseMsg(id)) }()

	if m, err := conn.Recv(); err != nil {
		return nil, errors.WithMessage(err, "Failed to receive message")
	} else if addrM, ok := m.(*AuthResponseMsg); !ok {
		return nil, errors.WithMessagef(err, "Expected AuthResponse wire msg, got %t", m.Type())
	} else {
		err := <-sent // Wait until the message was sent.
		return addrM.Address, err
	}
}

var _ msg.Msg = (*AuthResponseMsg)(nil)

// AuthResponseMsg is a message used to authenticate a peer with another peer.
type AuthResponseMsg struct {
	Address Address
}

func (m *AuthResponseMsg) Type() msg.Type {
	return msg.AuthResponse
}

func (m *AuthResponseMsg) Encode(w io.Writer) error {
	return m.Address.Encode(w)
}

func (m *AuthResponseMsg) Decode(r io.Reader) (err error) {
	m.Address, err = wallet.DecodeAddress(r)
	return err
}

// NewAuthResponseMsg creates an authentication response message.
// In the future, it will also take an authentication challenge message as
// additional argument.
func NewAuthResponseMsg(id Identity) msg.Msg {
	return &AuthResponseMsg{id.Address()}
}
