package fabapi

import (
	"github.com/getlantern/deepcopy"
	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/pkg/errors"
	"github.com/tuyendev/fabric-config/internal/configtx"
	"github.com/tuyendev/fabric-config/pkg/archetype"
)

func CalculateOrdererConfigAfterAdd(orderer archetype.Orderer, addingOrg archetype.OrdererOrganization, latestConfig *cb.ConfigGroup) ([]byte, error) {
	newOrderer := calculateNewOrdererAfterAdd(orderer, addingOrg)
	if newOrderer == nil {
		return nil, errors.Errorf("cannot calculate new orderer after add new order")
	}
	newOrderingProfile := getOrderingProfile(*newOrderer, getStandardOrderingCapabilities(), getStandardOrderingPolicies())
	newConfig, err := configtx.NewSystemChannelGroup(newOrderingProfile)
	if err != nil {
		return nil, err
	}
	updateBytes, err := calculateOrdererConfigUpdateAfterAdd(newConfig, latestConfig, *newOrderer)
	if err != nil {
		return nil, err
	}
	return CreateSignedConfigEnvelope(orderer.GetPlainOrganizations(), updateBytes)
}

func calculateNewOrdererAfterAdd(orderer archetype.Orderer, addingOrg archetype.OrdererOrganization) *archetype.Orderer {
	newOrderer := archetype.Orderer{}
	if err := deepcopy.Copy(&newOrderer, &orderer); err != nil {
		return nil
	}
	checkNewOrderAdded := false
	for i := range newOrderer.Organizations {
		org := &newOrderer.Organizations[i]
		if org.Name == addingOrg.Name {
			org.Orders = append(org.Orders, addingOrg.Orders...)
			checkNewOrderAdded = true
		}
	}
	if !checkNewOrderAdded {
		newOrderer.Organizations = append(newOrderer.Organizations)
	}
	return &newOrderer
}

//TODO Consider to update address only
func calculateOrdererConfigUpdateAfterAdd(newConfig, latestConfig *cb.ConfigGroup, newOrderer archetype.Orderer) ([]byte, error) {
	proposalConfigTx := configtx.New(&cb.Config{ChannelGroup: latestConfig})
	updatingConfig := proto.Clone(latestConfig).(*cb.ConfigGroup)

	//update consensus
	consensusType := updatingConfig.Groups[channelconfig.OrdererGroupKey].Values[channelconfig.ConsensusTypeKey]
	consensusType.Version = consensusType.Version + 1
	consensusType.Value = newConfig.Groups[channelconfig.OrdererGroupKey].Values[channelconfig.ConsensusTypeKey].Value

	isOrgUpdated := false
	for _, org := range newOrderer.Organizations {
		orgInNew, existedInNew := newConfig.Groups[channelconfig.OrdererGroupKey].Groups[org.GetCfgTxName()]
		if !existedInNew {
			return nil, errors.Errorf("org with name '%s' does not exist in config", org.Name)
		}
		_, existedInOld := updatingConfig.Groups[channelconfig.OrdererGroupKey].Groups[org.GetCfgTxName()]
		if !existedInOld {
			isOrgUpdated = true
			updatingConfig.Groups[channelconfig.OrdererGroupKey].Groups[org.GetCfgTxName()] = orgInNew
		}
	}
	if isOrgUpdated {
		updatingConfig.Groups[channelconfig.OrdererGroupKey].Version = updatingConfig.Groups[channelconfig.OrdererGroupKey].Version + 1
	}
	updatingConfig.Groups[channelconfig.OrdererGroupKey].Version = updatingConfig.Groups[channelconfig.OrdererGroupKey].Version + 1
	*proposalConfigTx.UpdatedConfig() = cb.Config{ChannelGroup: updatingConfig}
	return proposalConfigTx.ComputeMarshaledUpdate(newOrderer.Name)
}

func CalculateOrdererConfigAfterRemove(orderer archetype.Orderer, removingOrg archetype.OrdererOrganization, latestConfig *cb.ConfigGroup) ([]byte, error) {
	newOrderer := calculateNewOrdererAfterRemove(orderer, removingOrg)
	if newOrderer == nil {
		return nil, errors.Errorf("cannot calculate new orderer after add new order")
	}
	newOrderingProfile := getOrderingProfile(*newOrderer, getStandardOrderingCapabilities(), getStandardOrderingPolicies())
	newConfig, err := configtx.NewSystemChannelGroup(newOrderingProfile)
	if err != nil {
		return nil, err
	}
	updateBytes, err := calculateOrdererConfigUpdateAfterRemove(newConfig, latestConfig, orderer)
	if err != nil {
		return nil, err
	}
	return CreateSignedConfigEnvelope(orderer.GetPlainOrganizations(), updateBytes)
}

func calculateNewOrdererAfterRemove(orderer archetype.Orderer, removingOrg archetype.OrdererOrganization) *archetype.Orderer {
	newOrderer := archetype.Orderer{}
	if err := deepcopy.Copy(&newOrderer, &orderer); err != nil {
		return nil
	}
	orgMap := make(map[string]*archetype.OrdererOrganization)
	for i := range newOrderer.Organizations {
		org := newOrderer.Organizations[i]
		orgMap[org.GetCfgTxName()] = &org
	}
	potentialOrg := orgMap[removingOrg.GetCfgTxName()]
	var newOrders []archetype.Order
	removingOrder := removingOrg.Orders[0]
	for _, order := range potentialOrg.Orders {
		if order.Host == removingOrder.Host && order.ListenPort == removingOrder.ListenPort {
			continue
		}
		newOrders = append(newOrders, order)
	}
	if len(newOrders) == 0 {
		delete(orgMap, potentialOrg.GetCfgTxName())
	}
	potentialOrg.Orders = newOrders
	var newOrganizations []archetype.OrdererOrganization
	for _, organization := range orgMap {
		newOrganizations = append(newOrganizations, *organization)
	}
	newOrderer.Organizations = newOrganizations
	return &newOrderer
}

//TODO Consider to update address only
func calculateOrdererConfigUpdateAfterRemove(newConfig, latestConfig *cb.ConfigGroup, orderer archetype.Orderer) ([]byte, error) {
	proposalConfigTx := configtx.New(&cb.Config{ChannelGroup: latestConfig})
	updatingConfig := proto.Clone(latestConfig).(*cb.ConfigGroup)

	//update consensus
	consensusType := updatingConfig.Groups[channelconfig.OrdererGroupKey].Values[channelconfig.ConsensusTypeKey]
	consensusType.Version = consensusType.Version + 1
	consensusType.Value = newConfig.Groups[channelconfig.OrdererGroupKey].Values[channelconfig.ConsensusTypeKey].Value

	isOrgUpdated := false
	for _, org := range orderer.Organizations {
		_, existedInOld := updatingConfig.Groups[channelconfig.OrdererGroupKey].Groups[org.GetCfgTxName()]
		if !existedInOld {
			return nil, errors.Errorf("org with name '%s' does not exist in config", org.Organization.Name)
		}
		_, existedInNew := newConfig.Groups[channelconfig.OrdererGroupKey].Groups[org.GetCfgTxName()]
		if !existedInNew {
			isOrgUpdated = true
			delete(updatingConfig.Groups[channelconfig.OrdererGroupKey].Groups, org.GetCfgTxName())
		}
	}
	if isOrgUpdated {
		updatingConfig.Groups[channelconfig.OrdererGroupKey].Version = updatingConfig.Groups[channelconfig.OrdererGroupKey].Version + 1
	}
	updatingConfig.Groups[channelconfig.OrdererGroupKey].Version = updatingConfig.Groups[channelconfig.OrdererGroupKey].Version + 1
	*proposalConfigTx.UpdatedConfig() = cb.Config{ChannelGroup: updatingConfig}
	return proposalConfigTx.ComputeMarshaledUpdate(orderer.Name)
}
