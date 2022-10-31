# go-starknet

`go-starknet` is a CLI written in Go for Starknet. To install it, you can
simply:

## if you have Go 1.18+ installed

```shell
cd
go install github.com/dontpanicdao/caigo/go-starknet@latest
go-starknet help
```

## on MacOS with Homebrew

```shell
cd
brew tap dontpanicdao/dontpanicdao
brew install go-starknet
go-starknet help
```


usage: starknet [-h] [-v] [--network NETWORK] [--network_id NETWORK_ID] [--chain_id CHAIN_ID] [--wallet WALLET | --no_wallet]
                [--account ACCOUNT] [--account_dir ACCOUNT_DIR] [--flavor {Debug,Release,RelWithDebInfo}] [--show_trace]
                [--gateway_url GATEWAY_URL] [--feeder_gateway_url FEEDER_GATEWAY_URL]
                {call,declare,deploy,deploy_account,estimate_message_fee,get_block,get_block_traces,get_class_by_hash,get_class_hash_at,get_code,get_contract_addresses,get_full_contract,get_nonce,get_state_update,get_storage_at,get_transaction,get_transaction_receipt,get_transaction_trace,invoke,new_account,tx_status}
starknet: error: argument command: invalid choice: 'block' (choose from 'call', 'declare', 'deploy', 'deploy_account', 'estimate_message_fee', 'get_block', 'get_block_traces', 'get_class_by_hash', 'get_class_hash_at', 'get_code', 'get_contract_addresses', 'get_full_contract', 'get_nonce', 'get_state_update', 'get_storage_at', 'get_transaction', 'get_transaction_receipt', 'get_transaction_trace', 'invoke', 'new_account', 'tx_status')