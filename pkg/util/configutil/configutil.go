package configutil

import (
	"github.com/hyperledger/fabric-protos-go/common"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"
	"io/ioutil"
)

func GetConfigFromFile(filepath string) *cb.Config {
	latestConfig, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	block, err := protoutil.UnmarshalBlock(latestConfig)
	if err != nil {
		panic(err)
	}
	return calculateChannelConfigFromBlock(block)
}

func calculateChannelConfigFromBlock(configBlock *cb.Block) *cb.Config {
	envelopeConfig, err := protoutil.ExtractEnvelope(configBlock, 0)
	if err != nil {
		panic(err)
	}
	configEnv := &common.ConfigEnvelope{}
	_, err = protoutil.UnmarshalEnvelopeOfType(envelopeConfig, common.HeaderType_CONFIG, configEnv)
	if err != nil {
		panic(err)
	}
	return configEnv.Config
}
