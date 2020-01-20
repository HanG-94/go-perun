// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by a MIT-style license that can be found in
// the LICENSE file.

package channel

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"perun.network/go-perun/backend/ethereum/bindings/adjudicator"
	"perun.network/go-perun/backend/ethereum/bindings/assets"
	"perun.network/go-perun/log"
)

const deployGasLimit = 6200000

// DeployETHAssetholder deploys a new ETHAssetHolder contract.
func DeployETHAssetholder(ctx context.Context, backend ContractBackend, adjudicatorAddr common.Address) (common.Address, error) {
	auth, err := backend.newTransactor(ctx, big.NewInt(0), deployGasLimit)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "could not create transactor")
	}
	addr, tx, _, err := assets.DeployAssetHolderETH(auth, backend, adjudicatorAddr)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "could not create transaction")
	}
	if err := execSuccessful(ctx, backend, tx); err != nil {
		return common.Address{}, errors.WithMessage(err, "deploying ethassetholder")
	}
	log.Infof("Successfully deployed AssetHolderETH at %v.", addr.Hex())
	return addr, nil
}

// DeployAdjudicator deploys a new Adjudicator contract.
func DeployAdjudicator(ctx context.Context, backend ContractBackend) (common.Address, error) {
	auth, err := backend.newTransactor(ctx, big.NewInt(0), deployGasLimit)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "could not create transactor")
	}
	addr, tx, _, err := adjudicator.DeployAdjudicator(auth, backend)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "could not create transaction")
	}
	if err = execSuccessful(ctx, backend, tx); err != nil {
		return common.Address{}, errors.WithMessage(err, "deploying adjudicator")
	}
	log.Infof("Successfully deployed Adjudicator at %v.", addr.Hex())
	return addr, nil
}

func execSuccessful(ctx context.Context, backend ContractBackend, tx *types.Transaction) error {
	receipt, err := bind.WaitMined(ctx, backend, tx)
	if err != nil {
		return errors.Wrap(err, "could not execute transaction")
	}
	if receipt.Status == types.ReceiptStatusFailed {
		return errors.New("transaction failed")
	}
	return nil
}
