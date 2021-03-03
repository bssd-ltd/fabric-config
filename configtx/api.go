package configtx

import (
	"fmt"
	cb "github.com/hyperledger/fabric-protos-go/common"
)

func NewSystemChannelGroup(channelConfig Channel) (*cb.ConfigGroup, error) {
	return newSystemChannelGroup(channelConfig)
}

func NewApplicationChannelGroup(channelConfig Channel) (*cb.ConfigGroup, error) {
	channelGroup, err := newChannelGroupWithOrderer(channelConfig)
	if err != nil {
		return nil, err
	}

	applicationGroup, err := newApplicationGroupWithAnchorPeers(channelConfig.Application)
	if err != nil {
		return nil, err
	}

	channelGroup.Groups[ApplicationGroupKey] = applicationGroup

	channelGroup.ModPolicy = AdminsPolicyKey

	return channelGroup, nil
}

func newApplicationGroupWithAnchorPeers(application Application) (*cb.ConfigGroup, error) {
	applicationGroup, err := newApplicationGroupTemplate(application)
	if err != nil {
		return nil, err
	}

	for _, org := range application.Organizations {
		applicationGroup.Groups[org.Name], err = newApplicationOrgConfigGroup(org)
		if err != nil {
			return nil, fmt.Errorf("org group '%s': %v", org.Name, err)
		}
	}

	return applicationGroup, nil
}
