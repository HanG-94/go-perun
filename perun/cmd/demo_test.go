// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by a MIT-style license that can be found in
// the LICENSE file.

package cmd_test

import (
	"fmt"
	"net"
	"regexp"
	"testing"
	"time"

	expect "github.com/google/goexpect"
	"github.com/stretchr/testify/require"
)

var (
	any     = regexp.MustCompile(".*")
	timeout = time.Second * 20
)

func getBalances() (string, error) {
	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	fmt.Fprintf(conn, "getbals")
	buff := make([]byte, 1024)
	_, err = conn.Read(buff)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func TestNodes(t *testing.T) {
	alice, _, err := expect.Spawn("go run perun/main.go demo --config perun/alice.yaml --log warn --test-api true", -1)
	require.NoError(t, err)
	defer alice.Close()

	bob, _, err := expect.Spawn("go run perun/main.go demo --config perun/bob.yaml --log warn --test-api true", -1)
	require.NoError(t, err)
	defer bob.Close()

	// Alice start
	alice.Expect(any, timeout)
	// Bob start
	bob.Expect(any, timeout)
	// Alice connect to Bob
	alice.Send("connect 0.0.0.0 0x05e71027e7d3bd6261de7634cf50F0e2142067C4 bob\r")
	// Alice open channel to Bob
	alice.Send("open bob 10 10\r")
	// Alice wait for channel
	alice.Expect(any, timeout)
	// Alice get balances
	b, err := getBalances()
	require.NoError(t, err)
	t.Log("Balances: ", b)
}
