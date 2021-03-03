package archetype

import (
	"github.com/hyperledger/fabric-config/configtx"
	orderertx "github.com/hyperledger/fabric-config/configtx/orderer"
	"github.com/hyperledger/fabric-config/pkg/util/certreader"
	"github.com/hyperledger/fabric-config/pkg/util/strex"
	"strconv"
	"time"
)

const ordererMaxChannels uint64 = 10

type Orderer struct {
	Name            string
	Organizations   []OrdererOrganization
	MaxChannels     uint64
	Policies        map[string]configtx.Policy
	BatchTimeout    time.Duration
	BatchSize       *orderertx.BatchSize
	Capabilities    []string
	EtcdRaftOptions *orderertx.EtcdRaftOptions
	Consortiums     []Consortia
}

func (c *Orderer) GetPlainOrganizations() []Organization {
	var organizations []Organization
	for _, org := range c.Organizations {
		organizations = append(organizations, *org.Organization)
	}
	return organizations
}

func (o *Orderer) GetCfgTxEtcdRaftOrderer() configtx.Orderer {
	if o.MaxChannels == 0 {
		o.MaxChannels = o.getCfgTxStandardMaxChannels()
	}
	if o.BatchTimeout == 0 {
		o.BatchTimeout = o.getCfgTxStandardBatchTimeout()
	}
	if o.BatchSize == nil {
		o.BatchSize = o.getCfgTxStandardBatchSize()
	}
	if len(o.Capabilities) == 0 {
		o.Capabilities = o.getCfgTxStandardCapabilities()
	}
	if len(o.Policies) == 0 {
		o.Policies = o.getCfgTxStandardPolicies()
	}
	return configtx.Orderer{
		OrdererType:   orderertx.ConsensusTypeEtcdRaft,
		BatchTimeout:  o.BatchTimeout,
		BatchSize:     *o.BatchSize,
		EtcdRaft:      o.getCfgTxEtcdRaftConfig(),
		Organizations: o.getCfgTxOrgs(),
		MaxChannels:   o.MaxChannels,
		Capabilities:  o.Capabilities,
		Policies:      o.Policies,
		State:         orderertx.ConsensusStateNormal,
	}
}

func (o *Orderer) getCfgTxOrgs() []configtx.Organization {
	var result []configtx.Organization
	for _, oo := range o.Organizations {
		result = append(result, oo.getCfgTxOrg())
	}
	return result
}

func (*Orderer) getCfgTxStandardMaxChannels() uint64 {
	return ordererMaxChannels
}

func (*Orderer) getCfgTxStandardBatchTimeout() time.Duration {
	return 10 * time.Second
}

func (*Orderer) getCfgTxStandardBatchSize() *orderertx.BatchSize {
	return &orderertx.BatchSize{
		MaxMessageCount:   10,
		AbsoluteMaxBytes:  10 * 1024 * 1024,
		PreferredMaxBytes: 512 * 1024,
	}
}

func (*Orderer) getCfgTxStandardCapabilities() []string {
	return []string{"V2_0"}
}

func (*Orderer) getCfgTxStandardPolicies() map[string]configtx.Policy {
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
		configtx.BlockValidationPolicyKey: {
			Type: configtx.ImplicitMetaPolicyType,
			Rule: "ANY Writers",
		},
	}
}

func (o *Orderer) getCfgTxEtcdRaftConfig() orderertx.EtcdRaft {
	if o.EtcdRaftOptions == nil {
		o.EtcdRaftOptions = o.getCfgTxStandardEtcdRaftOptions()
	}
	return orderertx.EtcdRaft{
		Consenters: o.getCfgTxConsenters(),
		Options:    *o.EtcdRaftOptions,
	}
}

func (*Orderer) getCfgTxStandardEtcdRaftOptions() *orderertx.EtcdRaftOptions {
	return &orderertx.EtcdRaftOptions{
		TickInterval:         "500ms",
		ElectionTick:         10,
		HeartbeatTick:        1,
		MaxInflightBlocks:    5,
		SnapshotIntervalSize: 16 * 1024 * 1024,
	}
}

func (o *Orderer) getCfgTxConsenters() []orderertx.Consenter {
	var result []orderertx.Consenter
	for _, oo := range o.Organizations {
		result = append(result, oo.getCfgTxConsenters()...)
	}
	return result
}

func (o *Orderer) GetCfgTxConsortiums() []configtx.Consortium {
	var result []configtx.Consortium
	for _, consortium := range o.Consortiums {
		result = append(result, consortium.getCfgTxCosortium())
	}
	return result
}

type OrdererOrganization struct {
	*Organization
	Orders   []Order
	Policies map[string]configtx.Policy
}

func (oo *OrdererOrganization) getCfgTxOrg() configtx.Organization {
	if len(oo.Policies) == 0 {
		oo.Policies = oo.getCfgTxStandardOrganizationPolicies()
	}
	return configtx.Organization{
		Name:             oo.GetCfgTxName(),
		Policies:         oo.Policies,
		MSP:              oo.getCfgTxMSP(),
		OrdererEndpoints: oo.getCfgTxOrdererEndpoints(),
	}
}

func (oo *OrdererOrganization) getCfgTxConsenters() []orderertx.Consenter {
	var result []orderertx.Consenter
	for _, order := range oo.Orders {
		result = append(result, order.getCfgTxConsenter())
	}
	return result
}

func (oo *OrdererOrganization) getCfgTxOrdererEndpoints() []string {
	var endpoints []string
	for _, order := range oo.Orders {
		endpoints = append(endpoints, order.getCfgTxEndpoint())
	}
	return endpoints
}

func (oo *OrdererOrganization) getCfgTxStandardOrganizationPolicies() map[string]configtx.Policy {
	return map[string]configtx.Policy{
		configtx.ReadersPolicyKey: {
			Type: configtx.SignaturePolicyType,
			Rule: "OR('" + oo.getMspID() + ".member')",
		},
		configtx.WritersPolicyKey: {
			Type: configtx.SignaturePolicyType,
			Rule: "OR('" + oo.getMspID() + ".member')",
		},
		configtx.AdminsPolicyKey: {
			Type: configtx.SignaturePolicyType,
			Rule: "OR('" + oo.getMspID() + ".admin')",
		},
	}
}

type Order struct {
	Host              string
	ListenPort        int
	ClientTLSCertPath string
	ServerTLSCertPath string
}

func (order *Order) getCfgTxEndpoint() string {
	return strex.Join(":", order.Host, strconv.Itoa(order.ListenPort))
}

func (order *Order) getCfgTxConsenter() orderertx.Consenter {
	clientTLSCert := certreader.ReadX509Cert(order.ClientTLSCertPath)
	serverTLSCert := certreader.ReadX509Cert(order.ServerTLSCertPath)
	return orderertx.Consenter{
		Address: orderertx.EtcdAddress{
			Host: order.Host,
			Port: order.ListenPort,
		},
		ClientTLSCert: clientTLSCert,
		ServerTLSCert: serverTLSCert,
	}
}
