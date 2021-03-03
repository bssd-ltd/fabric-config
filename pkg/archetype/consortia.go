package archetype

import (
	"github.com/hyperledger/fabric-config/configtx"
)

type ConsortiaOrganization struct {
	*Organization
	Polices map[string]configtx.Policy
}

func (co *ConsortiaOrganization) getCfgTxOrg() configtx.Organization {
	if len(co.Polices) == 0 {
		co.Polices = co.getCfgTxStandardPolicies()
	}
	return configtx.Organization{
		Name:     co.GetCfgTxName(),
		Policies: co.Polices,
		MSP:      co.getCfgTxMSP(),
	}
}

func (oc *ConsortiaOrganization) getCfgTxStandardPolicies() map[string]configtx.Policy {
	mspId := oc.getMspID()
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

type Consortia struct {
	Name          string
	Organizations []ConsortiaOrganization
}

func (c *Consortia) getCfgTxOrganizations() []configtx.Organization {
	var result []configtx.Organization
	for _, org := range c.Organizations {
		result = append(result, org.getCfgTxOrg())
	}
	return result
}

func (c *Consortia) getCfgTxCosortium() configtx.Consortium {
	return configtx.Consortium{
		Name:          c.Name,
		Organizations: c.getCfgTxOrganizations(),
	}
}
