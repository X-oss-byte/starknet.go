package caigo

import (
	"context"
	"math/big"

	"github.com/dontpanicdao/caigo/types"
)

const (
	FEE_MARGIN uint64 = 115
)

var (
	TRANSACTION_VERSION = types.StrToFelt("0")
	TRANSACTION_PREFIX  = types.StrToFelt("invoke")
	EXECUTE_SELECTOR    = types.StrToFelt("__execute__")
)

type Account struct {
	Provider types.Provider
	Address  *types.Felt
	PublicX  *big.Int
	PublicY  *big.Int
	private  *big.Int
}

type ExecuteDetails struct {
	MaxFee  *types.Felt
	Nonce   *types.Felt
	Version *types.Felt // not used currently
}

/*
Instantiate a new StarkNet Account which includes structures for calling the network and signing transactions:
- private signing key
- stark curve definition
- full provider definition
- public key pair for signature verifications
*/
func NewAccount(private, address *types.Felt, provider types.Provider) (*Account, error) {
	x, y, err := Curve.PrivateToPoint(private)
	if err != nil {
		return nil, err
	}

	return &Account{
		Provider: provider,
		Address:  address,
		PublicX:  x,
		PublicY:  y,
		private:  private,
	}, nil
}

func (account *Account) Sign(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return Curve.Sign(msgHash, account.private)
}

/*
invocation wrapper for StarkNet account calls to '__execute__' contact calls through an account abstraction
- implementation has been tested against OpenZeppelin Account contract as of: https://github.com/OpenZeppelin/cairo-contracts/blob/4116c1ecbed9f821a2aa714c993a35c1682c946e/src/openzeppelin/account/Account.cairo
- accepts a multicall
*/
func (account *Account) Execute(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.AddTxResponse, error) {
	if details.Nonce == nil {
		nonce, err := account.Provider.AccountNonce(ctx, account.Address)
		if err != nil {
			return nil, err
		}
		details.Nonce = nonce
	}

	if details.MaxFee == nil {
		fee, err := account.EstimateFee(ctx, calls, details)
		if err != nil {
			return nil, err
		}
		details.MaxFee = types.SetUint64((fee.OverallFee * FEE_MARGIN) / 100)
	}

	req, err := account.fmtExecute(ctx, calls, details)
	if err != nil {
		return nil, err
	}

	return account.Provider.Invoke(ctx, *req)
}

func (account *Account) HashMultiCall(fee *types.Felt, nonce *types.Felt, calls []types.Transaction) (*big.Int, error) {
	chainID, err := account.Provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	callArray := ExecuteCalldata(nonce, calls)

	// convert callArray into a BigInt array
	callArrayBigInt := make([]*types.Felt, 0)
	for _, call := range callArray {
		callArrayBigInt = append(callArrayBigInt, call)
	}

	cdHash, err := Curve.ComputeHashOnElements(callArrayBigInt)
	if err != nil {
		return nil, err
	}

	multiHashData := []*types.Felt{
		TRANSACTION_PREFIX,
		TRANSACTION_VERSION,
		account.Address,
		GetSelectorFromName(EXECUTE_SELECTOR),
		cdHash,
		fee,
		types.StrToFelt(chainID),
	}

	return Curve.ComputeHashOnElements(multiHashData)
}

func (account *Account) EstimateFee(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.FeeEstimate, error) {
	if details.Nonce == nil {
		nonce, err := account.Provider.AccountNonce(ctx, account.Address)
		if err != nil {
			return nil, err
		}
		details.Nonce = nonce
	}

	if details.MaxFee == nil {
		details.MaxFee = &types.Felt{Int: big.NewInt(0)}
	}

	req, err := account.fmtExecute(ctx, calls, details)
	if err != nil {
		return nil, err
	}

	return account.Provider.EstimateFee(ctx, *req, "")
}

func (account *Account) fmtExecute(ctx context.Context, calls []types.Transaction, details ExecuteDetails) (*types.FunctionInvoke, error) {
	req := types.FunctionInvoke{
		FunctionCall: types.FunctionCall{
			ContractAddress:    account.Address,
			EntryPointSelector: EXECUTE_SELECTOR,
			Calldata:           ExecuteCalldata(details.Nonce, calls),
		},
		MaxFee: details.MaxFee,
	}

	hash, err := account.HashMultiCall(details.MaxFee, details.Nonce, calls)
	if err != nil {
		return nil, err
	}

	r, s, err := account.Sign(hash)
	if err != nil {
		return nil, err
	}
	req.Signature = types.Signature{types.BigToFelt(r), types.BigToFelt(s)}

	return &req, nil
}

/*
Formats the multicall transactions in a format which can be signed and verified by the network and OpenZeppelin account contracts
*/
func ExecuteCalldata(nonce *types.Felt, calls []types.Transaction) (calldataArray []*types.Felt) {
	callArray := []*types.Felt{types.SetUint64(len(calls))}

	for _, tx := range calls {
		callArray = append(callArray, tx.ContractAddress, GetSelectorFromName(tx.EntryPointSelector))

		if len(tx.Calldata) == 0 {
			callArray = append(callArray, types.BigToFelt(big.NewInt(0)), types.BigToFelt(big.NewInt(0)))

			continue
		}

		callArray = append(callArray, types.SetUint64(len(calldataArray)), types.Felt.SetUint64(len(tx.Calldata)))
		for _, cd := range tx.Calldata {
			calldataArray = append(calldataArray, cd)
		}
	}

	callArray = append(callArray, types.SetUint64(len(calldataArray)))
	callArray = append(callArray, calldataArray...)
	callArray = append(callArray, nonce)
	return callArray
}
