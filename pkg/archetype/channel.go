package archetype

import (
	"github.com/hyperledger/fabric-config/configtx"
)

type Application struct {
	Name          string
	Consortium    string
	Organizations []ChannelOrganization
	Capabilities  []string
	Policies      map[string]configtx.Policy
}

func (c *Application) GetPlainOrganizations() []Organization {
	var organizations []Organization
	for _, org := range c.Organizations {
		organizations = append(organizations, *org.Organization)
	}
	return organizations
}

func (c *Application) GetCfgTxApplication() configtx.Application {
	if len(c.Capabilities) == 0 {
		c.Capabilities = c.getCfgTxAppCapabilities()
	}
	if len(c.Policies) == 0 {
		c.Policies = c.getCfgTxAppPolicies()
	}
	return configtx.Application{
		Organizations: c.getCfgTxOrgs(),
		Capabilities:  c.Capabilities,
		Policies:      c.Policies,
	}
}

func (c *Application) getCfgTxOrgs() []configtx.Organization {
	var organizations []configtx.Organization
	for _, org := range c.Organizations {
		organizations = append(organizations, org.getCfgTxOrg())
	}
	return organizations
}

func (c *Application) getCfgTxAppCapabilities() []string {
	return []string{"V2_0"}
}

func (*Application) getCfgTxAppPolicies() map[string]configtx.Policy {
	return map[string]configtx.Policy{
		"Admins": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Admins",
		},
		"Endorsement": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Endorsement",
		},
		"LifecycleEndorsement": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Endorsement",
		},
		"Readers": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Readers",
		},
		"Writers": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Writers",
		},
	}
}

type ChannelOrganization struct {
	*Organization
	Peers    []Peer
	Policies map[string]configtx.Policy
}

func (co *ChannelOrganization) getCfgTxOrg() configtx.Organization {
	if len(co.Policies) == 0 {
		co.Policies = co.getCfgTxStandardPolicies()
	}
	return configtx.Organization{
		Name:        co.GetCfgTxName(),
		Policies:    co.Policies,
		MSP:         co.getCfgTxMSP(),
		AnchorPeers: co.getCfgTxAnchorPeers(),
	}
}

func (co *ChannelOrganization) getCfgTxStandardPolicies() map[string]configtx.Policy {
	mspId := co.getMspID()
	return map[string]configtx.Policy{
		configtx.ReadersPolicyKey: {
			Type: configtx.SignaturePolicyType,
			Rule: "OR('" + mspId + ".admin', '" + mspId + ".peer','" + mspId + ".client')",
		},
		configtx.WritersPolicyKey: {
			Type: configtx.SignaturePolicyType,
			Rule: "OR('" + mspId + ".admin', '" + mspId + ".client')",
		},
		configtx.AdminsPolicyKey: {
			Type: configtx.SignaturePolicyType,
			Rule: "OR('" + mspId + ".admin')",
		},
		configtx.EndorsementPolicyKey: {
			Type: configtx.SignaturePolicyType,
			Rule: "OR('" + mspId + ".peer')",
		},
	}
}

func (co *ChannelOrganization) getCfgTxAnchorPeers() []configtx.Address {
	var addresses []configtx.Address
	for _, peer := range co.Peers {
		addresses = append(addresses, peer.getCfgTxPeerAddress())
	}
	return addresses
}

type Peer struct {
	Host       string
	ListenPort int
}

func (peer *Peer) getCfgTxPeerAddress() configtx.Address {
	return configtx.Address{
		Host: peer.Host,
		Port: peer.ListenPort,
	}
}
