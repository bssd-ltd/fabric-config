package fabapi

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/tuyendev/fabric-config/internal/configtx"
	"github.com/tuyendev/fabric-config/pkg/archetype"
)

func CreateSignedConfigEnvelope(orgs []archetype.Organization, marshaledUpdate []byte) ([]byte, error) {
	configSignature, err := createConfigSignature(orgs, marshaledUpdate)
	if err != nil {
		return nil, err
	}
	envelope, err := configtx.NewEnvelope(marshaledUpdate, configSignature...)
	if err != nil {
		return nil, err
	}
	result, err := proto.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("marshaling config update: %v", err)
	}
	return result, nil
}

func createConfigSignature(orgs []archetype.Organization, marshaledUpdate []byte) ([]*common.ConfigSignature, error) {
	signingIdentities := getSigningIdentities(orgs)
	var result []*common.ConfigSignature
	for _, signer := range signingIdentities {
		cs, err := signer.CreateConfigSignature(marshaledUpdate)
		if err != nil {
			return nil, err
		}
		result = append(result, cs)
	}
	return result, nil
}

func getSigningIdentities(orgs []archetype.Organization) []*configtx.SigningIdentity {
	var result []*configtx.SigningIdentity
	for _, org := range orgs {
		result = append(result, org.GetSignIdentity())
	}
	return result
}
