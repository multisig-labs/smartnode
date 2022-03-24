package validator

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/smartnode/shared/types/eth2"
	eth2types "github.com/wealdtech/go-eth2-types/v2"

	"github.com/rocket-pool/smartnode/shared/services/beacon"
)

// Deposit settings
const DepositAmount = 16000000000 // gwei

// Get deposit data & root for a given validator key and withdrawal credentials
func GetDepositData(validatorKey *rsa.PrivateKey, withdrawalCredentials common.Hash, eth2Config beacon.Eth2Config) (eth2.DepositData, common.Hash, error) {

	// Build deposit data
	dd := eth2.DepositDataNoSignature{
		PublicKey:             x509.MarshalPKCS1PublicKey(&validatorKey.PublicKey)[:48],
		WithdrawalCredentials: withdrawalCredentials[:],
		Amount:                DepositAmount,
	}

	// Get signing root
	or, err := dd.HashTreeRoot()
	if err != nil {
		return eth2.DepositData{}, common.Hash{}, err
	}
	domain, err := eth2types.ComputeDomain(eth2types.DomainDeposit, eth2Config.GenesisForkVersion, eth2types.ZeroGenesisValidatorsRoot)
	if err != nil {
		return eth2.DepositData{}, common.Hash{}, err
	}
	sr := eth2.SigningRoot{
		ObjectRoot: or[:],
		Domain:     domain,
	}

	// Get signing root with domain
	srHash, err := sr.HashTreeRoot()
	if err != nil {
		return eth2.DepositData{}, common.Hash{}, err
	}

	sig, err := validatorKey.Sign(rand.Reader, srHash[:], crypto.SHA256)
	if err != nil {
		return eth2.DepositData{}, common.Hash{}, err
	}

	// Build deposit data struct (with signature)
	sig = sig[:96]
	var depositData = eth2.DepositData{
		PublicKey:             dd.PublicKey,
		WithdrawalCredentials: dd.WithdrawalCredentials,
		Amount:                dd.Amount,
		Signature:             sig,
	}

	// Get deposit data root
	depositDataRoot, err := depositData.HashTreeRoot()
	if err != nil {
		return eth2.DepositData{}, common.Hash{}, err
	}

	// Return
	return depositData, depositDataRoot, nil

}
