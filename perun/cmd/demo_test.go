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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	any     = regexp.MustCompile(".+")
	timeout = time.Second * 30
)

func getBalances() (string, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	fmt.Fprintf(conn, "getbals")
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return "", err
	}
	return string(buff[0:n]), nil
}

func sendSynchron(obj *expect.GExpect, str string) error {
	fmt.Printf("%s ", "$")
	for _, b := range []byte(str) {
		<-time.After(time.Millisecond * 10)
		if err := obj.Send(string([]byte{b})); err != nil {
			return err
		}
		fmt.Printf("%s", string([]byte{b}))
	}
	_, _, err := obj.Expect(any, timeout)
	return err
}

func TestNodes(t *testing.T) {
	alice, _, err := expect.Spawn("go run ../main.go demo --config ../alice.yaml --log warn --test-api true", -1)
	require.NoError(t, err)
	defer alice.Close()

	bob, _, err := expect.Spawn("go run ../main.go demo --config ../bob.yaml --log warn", -1)
	require.NoError(t, err)
	defer bob.Close()

	// Alice start
	_, _, e := alice.Expect(any, timeout)
	require.NoError(t, e)

	t.Log("Waiting for contract deployment")
	// Bob start
	_, _, e = bob.Expect(any, timeout)
	require.NoError(t, e)

	// Alice connect to Bob
	require.NoError(t, sendSynchron(alice, "connect 127.0.0.1 0x05e71027e7d3bd6261de7634cf50F0e2142067C4 bob\n"))
	t.Log("Alice connected")
	// Alice open channel to Bob
	require.NoError(t, sendSynchron(alice, "open bob 1000 1000\n"))
	t.Log("Opening channel…")
	<-time.After(time.Second * 5)
	// Alice send to Bob and Bob to Alice
	for i := 0; i < 25; i++ {
		require.NoError(t, sendSynchron(alice, "send bob 1\n"))
		require.NoError(t, sendSynchron(bob, "send peer0 2\n"))
	}

	// Alice get balances
	<-time.After(time.Second)
	b, err := getBalances()
	require.NoError(t, err)
	t.Log("Balances: ", b)
	assert.Equal(t, "{\"bob\":{\"My\":1025,\"Other\":975}}", b)
}
