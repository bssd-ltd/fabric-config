package fabapi

import (
	"github.com/golang/protobuf/proto"
	"github.com/tuyendev/fabric-config/internal/configtx"
	"github.com/tuyendev/fabric-config/pkg/archetype"
)

func GenerateStandardOrderingGenesisBlock(orderer archetype.Orderer) ([]byte, error) {
	return GenerateOrderingGenesisBlock(orderer, getStandardOrderingCapabilities(), getStandardOrderingPolicies())
}

func GenerateOrderingGenesisBlock(orderer archetype.Orderer, capabilities []string, policies map[string]configtx.Policy) ([]byte, error) {
	profile := getOrderingProfile(orderer, capabilities, policies)
	block, err := configtx.NewSystemChannelGenesisBlock(profile, orderer.Name)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(block)
}

func GenerateStandardApplicationGenesisBlock(application archetype.Application, orderer archetype.Orderer) ([]byte, error) {
	return GenerateApplicationGenesisBlock(application, orderer, getStandardApplicationCapabilities(), getStandardApplicationPolicies())
}

func GenerateApplicationGenesisBlock(application archetype.Application, orderer archetype.Orderer,
	capabilities []string, policies map[string]configtx.Policy) ([]byte, error) {
	profile := getApplicationProfile(application, orderer, capabilities, policies)
	configBlock, err := configtx.NewMarshaledCreateChannelTx(profile, application.Name)
	if err != nil {
		return nil, err
	}
	return CreateSignedConfigEnvelope(application.GetPlainOrganizations(), configBlock)
}
