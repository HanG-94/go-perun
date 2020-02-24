// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by a MIT-style license that can be found in
// the LICENSE file.

package main_test

import (
	//"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	//"perun.network/go-perun/log"

	expect "github.com/google/goexpect"
)

const (
	timeout = 8 * time.Second
)

var (
	promptRE = regexp.MustCompile(".*%.*")
)

func exp(t *testing.T, r *regexp.Regexp, e *expect.GExpect) {
	res, _, err := e.Expect(promptRE, timeout)
	if err != nil {
		t.Fatalf("Expected '%s' but got '%s'\n", r.String(), res)
	}
}

func TestDemo(t *testing.T) {
	e, _, err := expect.Spawn("go run main.go demo --config alice.yaml --log panic", timeout)
	require.NoError(t, err)
	defer e.Close()

	exp(t, promptRE, e)
}
