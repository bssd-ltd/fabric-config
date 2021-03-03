package fabapi

import (
	"github.com/tuyendev/fabric-config/internal/configtx"
	"github.com/tuyendev/fabric-config/pkg/archetype"
)

func getStandardOrderingCapabilities() []string {
	return []string{"V2_0"}
}

func getStandardOrderingPolicies() map[string]configtx.Policy {
	return map[string]configtx.Policy{
		configtx.ReadersPolicyKey: {
			Type: configtx.ImplicitMetaPolicyType,
			Rule: "ANY Readers",
		},
		configtx.WritersPolicyKey: {
			Type: configtx.ImplicitMetaPolicyType,
			Rule: "ANY Writers",
		},
		configtx.AdminsPolicyKey: {
			Type: configtx.ImplicitMetaPolicyType,
			Rule: "MAJORITY Admins",
		},
	}
}

func getOrderingProfile(orderer archetype.Orderer, capabilities []string, policies map[string]configtx.Policy) configtx.Channel {
	return configtx.Channel{
		Consortium:   orderer.Name,
		Orderer:      orderer.GetCfgTxEtcdRaftOrderer(),
		Consortiums:  orderer.GetCfgTxConsortiums(),
		Capabilities: capabilities,
		Policies:     policies,
	}
}

func getStandardApplicationCapabilities() []string {
	return []string{"V2_0"}
}

func getStandardApplicationPolicies() map[string]configtx.Policy {
	return map[string]configtx.Policy{
		configtx.ReadersPolicyKey: {
			Type: configtx.ImplicitMetaPolicyType,
			Rule: "ANY Readers",
		},
		configtx.WritersPolicyKey: {
			Type: configtx.ImplicitMetaPolicyType,
			Rule: "ANY Writers",
		},
		configtx.AdminsPolicyKey: {
			Type: configtx.ImplicitMetaPolicyType,
			Rule: "MAJORITY Admins",
		},
	}
}

func getApplicationProfile(application archetype.Application, orderer archetype.Orderer,
	capabilities []string, policies map[string]configtx.Policy) configtx.Channel {
	return configtx.Channel{
		Consortium:   application.Consortium,
		Orderer:      orderer.GetCfgTxEtcdRaftOrderer(),
		Application:  application.GetCfgTxApplication(),
		Capabilities: capabilities,
		Policies:     policies,
	}
}
