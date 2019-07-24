// Copyright (c) 2019 The Perun Authors. All rights reserved.
// This file is part of go-perun. Use of this source code is governed by a
// MIT-style license that can be found in the LICENSE file.

package log

import "os"

// None is a logger that doesn't do anything
var None Logger = &none{}

type none struct{}

func (none) Printf(string, ...interface{}) {}
func (none) Print(...interface{})          {}
func (none) Println(...interface{})        {}
func (none) Tracef(string, ...interface{}) {}
func (none) Debugf(string, ...interface{}) {}
func (none) Infof(string, ...interface{})  {}
func (none) Warnf(string, ...interface{})  {}
func (none) Errorf(string, ...interface{}) {}
func (none) Trace(...interface{})          {}
func (none) Debug(...interface{})          {}
func (none) Info(...interface{})           {}
func (none) Warn(...interface{})           {}
func (none) Error(...interface{})          {}
func (none) Traceln(...interface{})        {}
func (none) Debugln(...interface{})        {}
func (none) Infoln(...interface{})         {}
func (none) Warnln(...interface{})         {}
func (none) Errorln(...interface{})        {}

func (none) Fatal(...interface{})                      { os.Exit(1) }
func (none) Fatalf(string, ...interface{})             { os.Exit(1) }
func (none) Fatalln(...interface{})                    { os.Exit(1) }
func (none) Panic(args ...interface{})                 { panic(args) }
func (none) Panicf(format string, args ...interface{}) { panic(args) }
func (none) Panicln(args ...interface{})               { panic(args) }

func (n *none) WithField(key string, value interface{}) Logger {
	return n
}

func (n *none) WithFields(Fields) Logger {
	return n
}

func (n *none) WithError(error) Logger {
	return n
}
