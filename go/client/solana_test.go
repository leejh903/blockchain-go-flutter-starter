package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	rpcClient := rpc.New(rpc.DevNet_RPC)
	wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
	assert.NoError(t, err)

	ctx := context.TODO()
	recent, err := rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	assert.NoError(t, err)

	to := solana.MustPublicKeyFromBase58("9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6b")
	fmt.Println(to)
	from := solana.NewWallet()
	fmt.Println(from.PublicKey())

	// Airdrop 5 SOL to the new account:
	out, err := rpcClient.RequestAirdrop(
		ctx,
		from.PublicKey(),
		solana.LAMPORTS_PER_SOL*5,
		rpc.CommitmentFinalized,
	)
	fmt.Println(out)
	assert.NoError(t, err)

	time.Sleep(20 * time.Second)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				solana.LAMPORTS_PER_SOL,
				from.PublicKey(),
				to,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(from.PublicKey()),
	)
	assert.NoError(t, err)

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if from.PublicKey().Equals(key) {
				return &from.PrivateKey
			}
			return nil
		},
	)
	assert.NoError(t, err)

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		rpcClient,
		wsClient,
		tx,
	)
	fmt.Println(sig)
	assert.NoError(t, err)
}
