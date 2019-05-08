// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcr

import (
	"context"
	"fmt"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/devwarrior777/atomicswap/libs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type wallet struct {
	conn   *grpc.ClientConn
	client walletrpc.WalletServiceClient
}

// starRPC - starts a new GRPC client for the network and address specified
//            along with the certs path, in RPCInfo
func startRPC(testnet bool, rpcinfo libs.RPCInfo) (*wallet, error) {
	hostport, err := getNormalizedAddress(testnet, rpcinfo.HostPort)
	if err != nil {
		return nil, fmt.Errorf("wallet server address: %v", err)
	}
	creds, err := credentials.NewClientTLSFromFile(rpcinfo.Certs, "")
	if err != nil {
		return nil, fmt.Errorf("open certificate: %v", err)
	}
	wallet := &wallet{}
	// get a connection to the server
	wallet.conn, err = grpc.Dial(hostport, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %v", err)
	}
	// get a client
	wallet.client = walletrpc.NewWalletServiceClient(wallet.conn)

	return wallet, err
}

// stopRPC - Explicit stop when not using defer()
func (w *wallet) stopRPC() {
	w.conn.Close()
}

////////////////
// GRPC funcs //
////////////////

func (w *wallet) ping() error {
	request := &walletrpc.PingRequest{}
	_, err := w.client.Ping(context.Background(), request)
	return err
}

// // getBlockCount calls the getblockcount JSON-RPC method. It is
// // currently used as a simple 'ping' to discover if node RPC is available
// func getBlockCount(rpcclient *rpc.Client) (int, error) {
// 	rawResp, err := rpcclient.RawRequest("getblockcount", nil)
// 	if err != nil {
// 		return -1, err
// 	}
// 	var blockCount int
// 	err = json.Unmarshal(rawResp, &blockCount)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return blockCount, nil
// }

// func getTransaction(rpcclient *rpc.Client, txid string) (*libs.GetTxResult, error) {
// 	txidBytes, err := json.Marshal(txid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	param := []json.RawMessage{txidBytes}
// 	rawResp, err := rpcclient.RawRequest("gettransaction", param)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var resp struct {
// 		Confirmations uint64 `json:"confirmations"`
// 		Blockhash     string `json:"blockhash"`
// 		Blockindex    int    `json:"blockindex"`
// 		Blocktime     uint64 `json:"blocktime"`
// 		Time          uint64 `json:"time"`
// 		TimeReceived  uint64 `json:"timereceived"`
// 		Hex           string `json:"hex"`
// 	}

// 	err = json.Unmarshal(rawResp, &resp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result libs.GetTxResult
// 	result.Confirmations = resp.Confirmations
// 	result.Blockhash = resp.Blockhash
// 	result.Blockindex = resp.Blockindex
// 	result.Blocktime = resp.Blocktime
// 	result.Time = resp.Time
// 	result.TimeReceived = resp.TimeReceived
// 	result.Hex = resp.Hex
// 	return &result, nil
// }

// // getNewAddress calls the getnewaddress JSON-RPC method.
// func getNewAddress(testnet bool, rpcclient *rpc.Client) (xzcutil.Address, error) {
// 	chainParams := getChainParams(testnet)
// 	rawResp, err := rpcclient.RawRequest("getnewaddress", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var addrStr string
// 	err = json.Unmarshal(rawResp, &addrStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	addr, err := xzcutil.DecodeAddress(addrStr, chainParams)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !addr.IsForNet(chainParams) {
// 		return nil, fmt.Errorf("address %v is not intended for use on %v",
// 			addrStr, chainParams.Name)
// 	}
// 	if _, ok := addr.(*xzcutil.AddressPubKeyHash); !ok {
// 		return nil, fmt.Errorf("getnewaddress: address %v is not P2PKH",
// 			addr)
// 	}
// 	return addr, nil
// }

// // getRawChangeAddress calls the getrawchangeaddress JSON-RPC method.  It is
// // implemented manually as the rpcclient implementation always passes the
// // account parameter which was removed in Bitcoin Core 0.15.
// func getRawChangeAddress(testnet bool, rpcclient *rpc.Client) (xzcutil.Address, error) {
// 	chainParams := getChainParams(testnet)
// 	rawResp, err := rpcclient.RawRequest("getrawchangeaddress", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var addrStr string
// 	err = json.Unmarshal(rawResp, &addrStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	addr, err := xzcutil.DecodeAddress(addrStr, chainParams)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !addr.IsForNet(chainParams) {
// 		return nil, fmt.Errorf("address %v is not intended for use on %v",
// 			addrStr, chainParams.Name)
// 	}
// 	if _, ok := addr.(*xzcutil.AddressPubKeyHash); !ok {
// 		return nil, fmt.Errorf("getrawchangeaddress: address %v is not P2PKH",
// 			addr)
// 	}
// 	return addr, nil
// }

// // getFeePerKb queries the wallet for the transaction relay fee/kB to use and
// // the minimum mempool relay fee.  It first tries to get the user-set fee in the
// // wallet.  If unset, it attempts to find an estimate using estimatesmartfee 6.
// // If both of these fail, it falls back to mempool relay fee policy.
// //
// // For Zcoin this will always fall back until there is a statistically significant
// // number of transactions per block
// func getFeePerKb(rpcclient *rpc.Client) (useFee, relayFee xzcutil.Amount, err error) {
// 	var estimateResp struct {
// 		FeeRate float64 `json:"feerate"`
// 	}
// 	info, err := rpcclient.GetInfo()
// 	if err != nil {
// 		return 0, 0, fmt.Errorf("getinfo: %v", err)
// 	}
// 	relayFee, err = xzcutil.NewAmount(info.RelayFee)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	maxFee := info.PaytxFee
// 	if info.PaytxFee != 0 {
// 		if info.RelayFee > maxFee {
// 			maxFee = info.RelayFee
// 		}
// 		useFee, err = xzcutil.NewAmount(maxFee)
// 		return useFee, relayFee, err
// 	}

// 	params := []json.RawMessage{[]byte("6")}
// 	estimateRawResp, err := rpcclient.RawRequest("estimatesmartfee", params)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	err = json.Unmarshal(estimateRawResp, &estimateResp)
// 	if err == nil && estimateResp.FeeRate > 0 {
// 		useFee, err = xzcutil.NewAmount(estimateResp.FeeRate)
// 		if relayFee > useFee {
// 			useFee = relayFee
// 		}
// 		return useFee, relayFee, err
// 	}

// 	fmt.Println("warning: falling back to mempool relay fee policy")
// 	useFee, err = xzcutil.NewAmount(info.RelayFee)
// 	return useFee, relayFee, err
// }

// // fundRawTransaction calls the fundrawtransaction JSON-RPC method.  It is
// // implemented manually as client support is currently missing from the
// // xzcd/rpcclient package.
// func fundRawTransaction(rpcclient *rpc.Client, tx *wire.MsgTx, feePerKb xzcutil.Amount) (fundedTx *wire.MsgTx, fee xzcutil.Amount, err error) {
// 	var buf bytes.Buffer
// 	buf.Grow(tx.SerializeSize())
// 	tx.Serialize(&buf)
// 	param0, err := json.Marshal(hex.EncodeToString(buf.Bytes()))
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	param1, err := json.Marshal(struct {
// 		FeeRate float64 `json:"feeRate"`
// 	}{
// 		FeeRate: feePerKb.ToXZC(),
// 	})
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	params := []json.RawMessage{param0, param1}
// 	rawResp, err := rpcclient.RawRequest("fundrawtransaction", params)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	var resp struct {
// 		Hex       string  `json:"hex"`
// 		Fee       float64 `json:"fee"`
// 		ChangePos float64 `json:"changepos"`
// 	}
// 	err = json.Unmarshal(rawResp, &resp)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	fundedTxBytes, err := hex.DecodeString(resp.Hex)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	fundedTx = &wire.MsgTx{}
// 	err = fundedTx.Deserialize(bytes.NewReader(fundedTxBytes))
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	feeAmount, err := xzcutil.NewAmount(resp.Fee)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return fundedTx, feeAmount, nil
// }

// // createSig creates and returns the serialized raw signature and compressed
// // pubkey for a transaction input signature.  Due to limitations of the Zcoin
// // Core RPC API, this requires dumping a private key and signing in the client,
// // rather than letting the wallet sign.
// func createSig(testnet bool, tx *wire.MsgTx, idx int, pkScript []byte, addr xzcutil.Address,
// 	rpcclient *rpc.Client) (sig, pubkey []byte, err error) {

// 	wif, err := dpk(testnet, rpcclient, addr)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	sig, err = txscript.RawTxInSignature(tx, idx, pkScript, txscript.SigHashAll, wif.PrivKey)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return sig, wif.PrivKey.PubKey().SerializeCompressed(), nil
// }

// func dpk(testnet bool, rpcclient *rpc.Client, addr xzcutil.Address) (wif *xzcutil.WIF, err error) {
// 	chainParams := getChainParams(testnet)
// 	addrStr := addr.EncodeAddress()
// 	if !addr.IsForNet(chainParams) {
// 		return nil, fmt.Errorf("address %v is not intended for use on %v",
// 			addrStr, chainParams.Name)
// 	}
// 	param0, err := json.Marshal(addrStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	params := []json.RawMessage{param0}
// 	// This should always fail the first time as Zcoin added a one-time authoriz-
// 	// ation key returned in error string. Along with a warning. The idea is that
// 	// inexperienced people are warned if scammers propose they use `dumpprivkey'
// 	_, err = rpcclient.RawRequest("dumpprivkey", params)
// 	if err == nil {
// 		unexpected := errors.New("dpk: No authorization challenge")
// 		return nil, unexpected
// 	}

// 	errStr := err.Error()
// 	searchStr := "authorization code is: "
// 	i0 := strings.Index(errStr, searchStr)
// 	if i0 == -1 {
// 		return nil, err
// 	}
// 	i := i0 + len(searchStr)
// 	authStr := errStr[i : i+4]
// 	//
// 	param1, err := json.Marshal(authStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	params2 := []json.RawMessage{param0, param1}
// 	rawResp2, err := rpcclient.RawRequest("dumpprivkey", params2)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var sk string
// 	err = json.Unmarshal(rawResp2, &sk)
// 	if err != nil {
// 		return nil, err
// 	}

// 	w, err := xzcutil.DecodeWIF(sk)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return w, nil
// }

// func sendRawTransaction(rpcclient *rpc.Client, tx *wire.MsgTx) (*chainhash.Hash, error) {
// 	txHash, err := rpcclient.SendRawTransaction(tx, false)
// 	if err != nil {
// 		return nil, fmt.Errorf("sendrawtransaction: %v", err)
// 	}
// 	return txHash, nil
// }
