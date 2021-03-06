package ethereum

import (
	"math/big"

	"fin4-core/server/env"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lytics/logrus"
	"log"
)

// Ethereum ethereum struct to implement crypto interface
type Ethereum struct {
	rpc      *ethclient.Client
	sim      *backends.SimulatedBackend
	auth     *bind.TransactOpts
	keystore *keystore.KeyStore
}

// MustNewEthereum create new Ethereum interface, panic if no connection
func MustNewEthereum() *Ethereum {
	conn, err := ethclient.Dial("http://127.0.0.1:7545") //ethclient.Dial("https://rinkeby.infura.io/")

	if err != nil {
		logrus.Fatal("Failed to connect to the Ethereum client: %v", err)
		return nil
  }
	// server key
	rawKey := "f212fe8e6f935ee29dba7edc3da12af967fff2f8920505835cf2d0640da6c7fb"//env.MustGetenv("ETH_KEY_RAW")
	log.Println(rawKey)
	rawKeyECDSA, err := crypto.HexToECDSA(rawKey)
	if err != nil {
		logrus.Fatal("Something wrong with server private key.", err)
	}
	ks := keystore.NewKeyStore(
		env.MustGetenv("ETH_KEY_STORE_DIR"),
		keystore.LightScryptN,
		keystore.LightScryptP)
	ks.ImportECDSA(rawKeyECDSA, "passphrase")
	// Create an authorized transactor
	auth := bind.NewKeyedTransactor(rawKeyECDSA)
	if err != nil {
		logrus.Fatal("Failed to create transactor: %v", err)
	}
	// Setup blockchain simmulator
	gAlloc := map[common.Address]core.GenesisAccount{
		auth.From: {Balance: big.NewInt(10000000000)},
	}
	sim := backends.NewSimulatedBackend(gAlloc, 40000)
	return &Ethereum{
		rpc:      conn,
		sim:      sim,
		auth:     auth,
		keystore: ks,
	}
}

// CreateNewAddress returns best blocknumber in the blockchain
func (b *Ethereum) CreateNewAddress() (string, error) {
	acc, err := b.keystore.NewAccount("demo1")
	return acc.Address.String(), err
}

// DeployMintable deployes new Mintable token to Ethereum from server account
func (b *Ethereum) DeployMintable(
	name_ string,
	symbol_ string,
	decimals_ uint8,
	minter common.Address,
) (common.Address, *types.Transaction, error) {
	address, tx, _, err := DeployMintable(
		b.auth,
		// change here to rpc and it will deploy to rpc
		b.rpc,
		name_,
		symbol_,
		decimals_,
		minter,
	)
	if err != nil {
		return address, nil, err
	}
	return address, tx, nil
}

// Mint mints a new currency units to the passed token and toAddress
func (b *Ethereum) Mint(
	tokenAddress common.Address,
	toAddress common.Address,
	amount int64,
) (*types.Transaction, error) {
	// here change b.sim to rpc and it will communicate with the rpc
	mintable, err := NewMintable(tokenAddress, b.rpc)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("ethereum:Mint:1")
		return nil, err
	}
	// @TODO change b.auth to the address of the user who is minting the new
	// tokens i.e claim approver
	txAddress, err := mintable.Mint(b.auth, toAddress, big.NewInt(amount))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("ethereum:Mint:2")
		return nil, err
	}
	return txAddress, nil
}
