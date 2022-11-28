package secp256k1

import (
	"github.com/ElrondNetwork/wasm-vm-v1_2/crypto/signing"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type secp256k1 struct {
}

func NewSecp256k1() *secp256k1 {
	return &secp256k1{}
}

func (sec *secp256k1) VerifySecp256k1(key []byte, msg []byte, sig []byte) error {
	pubKey, err := btcec.ParsePubKey(key)
	if err != nil {
		return err
	}

	signature, err := ecdsa.ParseSignature(sig)
	if err != nil {
		return err
	}

	messageHash := chainhash.DoubleHashB(msg)
	verified := signature.Verify(messageHash, pubKey)

	if !verified {
		return signing.ErrInvalidSignature
	}

	return nil
}
