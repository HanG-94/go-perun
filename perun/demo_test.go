// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by a MIT-style license that can be found in
// the LICENSE file.

package main_test

import (
	"regexp"
	"testing"
	"time"

	expect "github.com/google/goexpect"
	"github.com/stretchr/testify/require"
)

const (
	timeout = 5 * time.Second
)

var (
	promptRE = regexp.MustCompile("^\\> ")
)

func TestDemo(t *testing.T) {
	e, _, err := expect.Spawn("date", timeout)
	require.NoError(t, err)
	defer e.Close()

	/*go func() {
		for {
			buff := make([]byte, 1024)
			_, err := e.Read(buff)
			if err != nil {
				log.Errorln("Could not read from program")
			}
			fmt.Printf("%s", buff)
		}
	}()*/

	_, _, err = e.Expect(promptRE, timeout)
	require.NoError(t, err)

}
