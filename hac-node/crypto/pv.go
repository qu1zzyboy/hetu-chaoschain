package crypto

import (
	"fmt"
	"github.com/cometbft/cometbft/crypto"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	cmtos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/privval"
	"os"
)

type PV struct {
	privateKey crypto.PrivKey
	publicKey  crypto.PubKey
}

func LoadFilePV(keyFilePath string) *PV {
	keyJSONBytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		cmtos.Exit(err.Error())
	}
	pvKey := privval.FilePVKey{}
	err = cmtjson.Unmarshal(keyJSONBytes, &pvKey)
	if err != nil {
		cmtos.Exit(fmt.Sprintf("Error reading PrivValidator key from %v: %v\n", keyFilePath, err))
	}

	return &PV{
		privateKey: pvKey.PrivKey,
		publicKey:  pvKey.PubKey,
	}
}

func (k *PV) PublicKey() []byte {
	return k.publicKey.Bytes()
}

func (k *PV) Address() string {
	return k.publicKey.Address().String()
}

func (k *PV) Sign(data []byte) ([]byte, error) {
	return k.privateKey.Sign(data)
}
