package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/urfave/cli/v2"
)

var classHashAtCommand = cli.Command{
	Name:   "get_class_hash_at",
	Usage:  "go-starknet get_class_hash_at 0x<contract address>",
	Flags:  classHashAtFlags,
	Action: classHashAtAction,
}

var classHashAtFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "provider",
		Usage: "choose between the gateway and rpc",
		Value: "gateway",
	},
	&cli.StringFlag{
		Name:  "base-url",
		Usage: "change the default baseURL",
		Value: "",
	},
}

func classHashAtAction(cCtx *cli.Context) error {
	providerName := cCtx.Value("provider")
	if providerName.(string) != "gateway" {
		return fmt.Errorf("provider not supported")
	}
	ctx := context.TODO()
	gw := gateway.NewProvider()
	contractAddress := cCtx.Args().First()
	if contractAddress == "" {
		return errors.New("contract address is mandatory")
	}
	contractAddressInt, ok := big.NewInt(0).SetString(contractAddress, 0)
	if !ok {
		return errors.New("address should start with 0x and be an hex")
	}
	hash, err := gw.ClassHashAt(ctx, fmt.Sprintf("0x%s", contractAddressInt.Text(16)))
	if err != nil {
		return err
	}
	fmt.Println(hash)
	return nil
}
