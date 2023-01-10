package factory

import (
	"github.com/multiversx/mx-chain-vm-go-v1_2/crypto"
	"github.com/multiversx/mx-chain-vm-go-v1_2/crypto/hashing"
	"github.com/multiversx/mx-chain-vm-go-v1_2/crypto/signing/bls"
	"github.com/multiversx/mx-chain-vm-go-v1_2/crypto/signing/ed25519"
	"github.com/multiversx/mx-chain-vm-go-v1_2/crypto/signing/secp256k1"
)

// NewVMCrypto returns a composite struct containing VMCrypto functionality implementations
func NewVMCrypto() crypto.VMCrypto {
	return struct {
		crypto.Hasher
		crypto.Ed25519
		crypto.BLS
		crypto.Secp256k1
	}{
		Hasher:    hashing.NewHasher(),
		Ed25519:   ed25519.NewEd25519Signer(),
		BLS:       bls.NewBLS(),
		Secp256k1: secp256k1.NewSecp256k1(),
	}
}
