package archetype

import (
	"crypto/x509"
	"github.com/hyperledger/fabric-config/configtx"
	"github.com/hyperledger/fabric-config/configtx/membership"
	"github.com/tuyendev/fabric-config/pkg/util/certreader"
	"github.com/tuyendev/fabric-config/pkg/util/strex"
	"path/filepath"
)

type Organization struct {
	Name   string
	MspDir string
	TlsDir string
}

func (org *Organization) GetCfgTxName() string {
	return org.getMspID()
}

func (org *Organization) getMspID() string {
	return strex.Concat(strex.TrimAndTitle(org.Name), "Msp")
}

func (org *Organization) getCfgTxMSP() configtx.MSP {
	ca := certreader.ReadX509Cert(filepath.Join(org.MspDir, RootCert))
	tlsca := certreader.ReadX509Cert(filepath.Join(org.TlsDir, RootCert))
	return org.getStandardCfgTxMSP(org.getMspID(), ca, tlsca)
}

func (org *Organization) getStandardCfgTxMSP(mspId string, ca *x509.Certificate, tlsca *x509.Certificate) configtx.MSP {
	return configtx.MSP{
		Name:      mspId,
		RootCerts: []*x509.Certificate{ca},
		CryptoConfig: membership.CryptoConfig{
			SignatureHashFamily:            "SHA2",
			IdentityIdentifierHashFunction: "SHA256",
		},
		TLSRootCerts: []*x509.Certificate{tlsca},
		NodeOUs: membership.NodeOUs{
			Enable: true,
			ClientOUIdentifier: membership.OUIdentifier{
				Certificate:                  ca,
				OrganizationalUnitIdentifier: MspTypeClient,
			},
			PeerOUIdentifier: membership.OUIdentifier{
				Certificate:                  ca,
				OrganizationalUnitIdentifier: MspTypePeer,
			},
			AdminOUIdentifier: membership.OUIdentifier{
				Certificate:                  ca,
				OrganizationalUnitIdentifier: MspTypeAdmin,
			},
			OrdererOUIdentifier: membership.OUIdentifier{
				Certificate:                  ca,
				OrganizationalUnitIdentifier: MspTypeOrderer,
			},
		},
	}
}

func (org *Organization) GetSignIdentity() *configtx.SigningIdentity {
	mspSignerCertPath := filepath.Join(org.MspDir, SignerMspCert)
	mspSignerKeyPath := filepath.Join(org.MspDir, SignerMspKey)
	return &configtx.SigningIdentity{
		Certificate: certreader.ReadX509Cert(mspSignerCertPath),
		PrivateKey:  certreader.ReadPrivateKey(mspSignerKeyPath),
		MSPID:       org.getMspID(),
	}
}
