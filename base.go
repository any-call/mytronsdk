package mytronsdk

import (
	"fmt"
	"github.com/any-call/mytronsdk/address"
	"github.com/any-call/mytronsdk/client"
	"github.com/any-call/mytronsdk/client/transaction"
	"github.com/any-call/mytronsdk/keys"
	"github.com/any-call/mytronsdk/proto/api"
	"math/big"
)

// 转账 trc20
func TransferTrc20(privateKeyHex string, toAddrT string, usdt float64) (string, error) {
	funderKey, err := keys.GetPrivateKeyFromHex(privateKeyHex)
	if err != nil {
		return "", err
	}

	funderAddr := address.BTCECPrivkeyToAddress(funderKey)

	c := client.NewGrpcClient("grpc.trongrid.io:50051")
	err = c.Start(client.GRPCInsecure())
	if err != nil {
		return "", err
	}
	defer c.Stop()

	toAddr, err := address.Base58ToAddress(toAddrT)
	if err != nil {
		return "", err
	}

	contractAddr := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" //合约地址
	amountBig := big.NewInt(int64(usdt * 1_000_000))     // USDT 6位精度
	method := "transfer(address,uint256)"
	data := fmt.Sprintf(`[{"address": "%s"}, {"uint256": %s}]`, toAddr.String(), amountBig.String())

	tx, err := c.TriggerContract(funderAddr.String(), contractAddr, //合约地址
		method, data, 20_000_000, 0, "", 0)
	if err != nil {
		fmt.Println("run here 4:", data)
		return "", err
	}

	signedTx, err := transaction.SignTransaction(tx.Transaction, funderKey)
	if err != nil {
		return "", err
	}

	result, err := c.Broadcast(signedTx)
	if err != nil {
		return "", err
	}

	if !result.Result ||
		result.Code != api.Return_SUCCESS {
		return "", fmt.Errorf("broadcast error: %s", result.Message)
	}

	return fmt.Sprintf("%x", tx.Txid), nil
}

// 转账 trx
func TransferTrx(privateKeyHex string, toAddrT string, trx float64) (string, error) {
	// Connect to network
	c := client.NewGrpcClient("grpc.trongrid.io:50051")
	err := c.Start(client.GRPCInsecure())
	if err != nil {
		return "", err
	}
	defer c.Stop()

	funderKey, err := keys.GetPrivateKeyFromHex(privateKeyHex)
	if err != nil {
		return "", err
	}

	funderAddr := address.BTCECPrivkeyToAddress(funderKey)

	tx, err := c.Transfer(funderAddr.String(), toAddrT, int64(trx*1_000_000)) // 100 TRX = 100 * 1e6 sun
	if err != nil {
		return "", err
	}

	signedTx, err := transaction.SignTransaction(tx.Transaction, funderKey)
	if err != nil {
		return "", err
	}

	result, err := c.Broadcast(signedTx)
	if err != nil {
		return "", err
	}

	if !result.Result ||
		result.Code != api.Return_SUCCESS {
		return "", fmt.Errorf("Broadcast failed: (%d) %s", result.Code, result.Message)
	}

	return fmt.Sprintf("%x", tx.Txid), nil
}
