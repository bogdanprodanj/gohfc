/*
Copyright: Cognition Foundry. All Rights Reserved.
License: Apache License Version 2.0
*/
package gohfc

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
)

// Identity is participant public and private key
type Identity struct {
	Certificate *x509.Certificate
	PrivateKey  interface{}
	MspId       string
}

// MarshalBinary encodes identity into a binary form and returns the result
func (i *Identity) MarshalBinary() (data []byte, err error) {
	var pk, cert string
	switch i.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		cast := i.PrivateKey.(*ecdsa.PrivateKey)
		b, err := x509.MarshalECPrivateKey(cast)
		if err != nil {
			return nil, err
		}
		block := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
		pk = base64.RawStdEncoding.EncodeToString(block)

	default:
		return nil, ErrInvalidKeyType
	}

	cert = base64.RawStdEncoding.EncodeToString(i.Certificate.Raw)
	return json.Marshal(map[string]string{"cert": cert, "pk": pk, "mspid": i.MspId})
}

// UnmarshalBinary unmarshal a binary representation of identity
func (i *Identity) UnmarshalBinary(data []byte) error {
	var raw map[string]string
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return err
	}
	// check do we have all keys
	if _, ok := raw["cert"]; !ok || len(raw["cert"]) < 1 {
		return ErrInvalidDataForParcelIdentity
	}
	if _, ok := raw["pk"]; !ok || len(raw["pk"]) < 1 {
		return ErrInvalidDataForParcelIdentity
	}

	certRaw, err := base64.RawStdEncoding.DecodeString(raw["cert"])
	if err != nil {
		return err
	}
	cert, err := x509.ParseCertificate(certRaw)
	if err != nil {
		return err
	}

	keyRaw, err := base64.RawStdEncoding.DecodeString(raw["pk"])
	if err != nil {
		return err
	}
	keyPem, _ := pem.Decode(keyRaw)
	if keyPem == nil {
		return ErrInvalidDataForParcelIdentity
	}
	var pk interface{}
	switch keyPem.Type {
	case "EC PRIVATE KEY":
		pk, err = x509.ParseECPrivateKey(keyPem.Bytes)
		if err != nil {
			return ErrInvalidDataForParcelIdentity
		}
	default:
		return ErrInvalidDataForParcelIdentity
	}

	i.Certificate = cert
	i.PrivateKey = pk
	i.MspId = raw["mspid"]

	return nil
}

// EnrollmentId get enrollment id from certificate
func (i *Identity) EnrollmentId() string {
	return i.Certificate.Subject.CommonName
}

// EnrollmentId get enrollment id from certificate
func (i *Identity) ToPem() ([]byte, []byte, error) {

	switch i.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		cast := i.PrivateKey.(*ecdsa.PrivateKey)
		b, err := x509.MarshalECPrivateKey(cast)
		if err != nil {
			return nil, nil, err
		}
		privateKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
		cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: i.Certificate.Raw})
		return cert, privateKey, nil

	default:
		return nil, nil, ErrInvalidKeyType
	}
}

// LoadCertFromFile read public key (pk) and private/secret kye (sk) from file system and return new Identity
func LoadCertFromFile(pk, sk string) (*Identity, error) {
	cf, err := ioutil.ReadFile(pk)
	if err != nil {
		return nil, err
	}

	kf, err := ioutil.ReadFile(sk)
	if err != nil {
		return nil, err
	}
	cpb, _ := pem.Decode(cf)
	kpb, _ := pem.Decode(kf)
	crt, err := x509.ParseCertificate(cpb.Bytes)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS8PrivateKey(kpb.Bytes)
	if err != nil {
		return nil, err
	}
	return &Identity{Certificate: crt, PrivateKey: key}, nil
}
