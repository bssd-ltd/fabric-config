package fabapi

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-config/configtx"
	"github.com/hyperledger/fabric-config/pkg/archetype"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/pkg/errors"
)

func UpdatePeer(application archetype.Application, orderer archetype.Orderer, latestConfig *cb.Config, actionType ActionType) ([]byte, error) {
	profile := getApplicationProfile(application, orderer, getStandardApplicationCapabilities(), getStandardApplicationPolicies())
	newConfig, err := configtx.NewApplicationChannelGroup(profile)
	if err != nil {
		return nil, err
	}
	proposalChannelConfigTx := configtx.New(&cb.Config{ChannelGroup: latestConfig.ChannelGroup})
	updatedConfig := proto.Clone(latestConfig.GetChannelGroup()).(*cb.ConfigGroup)
	updatedConfig.Groups[channelconfig.ApplicationGroupKey].Version = updatedConfig.Groups[channelconfig.ApplicationGroupKey].Version + 1
	switch actionType {
	case ActionAdd:
		for _, org := range application.Organizations {
			newer, existInNew := newConfig.Groups[channelconfig.ApplicationGroupKey].Groups[org.GetCfgTxName()]
			if !existInNew {
				return nil, errors.Errorf("org with name '%s' does not exist in config", org.Organization.Name)
			}
			updating, existInOld := updatedConfig.Groups[channelconfig.ApplicationGroupKey].Groups[org.GetCfgTxName()]
			if existInOld {
				updating.Values[channelconfig.AnchorPeersKey] = newer.Values[channelconfig.AnchorPeersKey]
			} else {
				updating = newer
			}
		}
		break
	case ActionRemove:
		for _, org := range application.Organizations {
			updating, existInOld := updatedConfig.Groups[channelconfig.ApplicationGroupKey].Groups[org.GetCfgTxName()]
			if !existInOld {
				return nil, errors.Errorf("org with name '%s' does not exist in config", org.Organization.Name)
			}
			newer, existInNew := newConfig.Groups[channelconfig.ApplicationGroupKey].Groups[org.GetCfgTxName()]
			if !existInNew {
				delete(updatedConfig.Groups[channelconfig.ApplicationGroupKey].Groups, org.GetCfgTxName())
			} else {
				updating.Values[channelconfig.AnchorPeersKey] = newer.Values[channelconfig.AnchorPeersKey]
			}
		}
		break
	case ActionUpdate:
		for _, org := range application.Organizations {
			updating, ok := updatedConfig.Groups[channelconfig.ApplicationGroupKey].Groups[org.GetCfgTxName()]
			if !ok {
				return nil, errors.Errorf("org with name '%s' does not exist in config", org.Organization.Name)
			}
			newest, ok := newConfig.Groups[channelconfig.ApplicationGroupKey].Groups[org.GetCfgTxName()]
			if !ok {
				return nil, errors.Errorf("org with name '%s' does not exist in config", org.Organization.Name)
			}
			updating.Values[channelconfig.AnchorPeersKey] = newest.Values[channelconfig.AnchorPeersKey]
		}
		break
	default:
		return nil, errors.Errorf("unkown action with peer")
	}
	*proposalChannelConfigTx.UpdatedConfig() = cb.Config{ChannelGroup: updatedConfig}
	marshaledUpdate, err := proposalChannelConfigTx.ComputeMarshaledUpdate(application.Name)
	if err != nil {
		return nil, err
	}
	return CreateSignedConfigEnvelope(application.GetPlainOrganizations(), marshaledUpdate)
}
