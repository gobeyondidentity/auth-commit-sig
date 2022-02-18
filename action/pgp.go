package action

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	// TODO: The OpenPGP package is deprecated. See
	// https://github.com/golang/go/issues/44226.
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// ParseSignatureIssuerKeyID parses an ASCII-armored PGP signature and extracts
// the PGP Key ID of the key that produced it.
func ParseSignatureIssuerKeyID(armoredSignature string) (string, error) {
	signature, err := parseSignature(armoredSignature)
	if err != nil {
		return "", fmt.Errorf("failed to parse signature: %w", err)
	}

	if signature.IssuerKeyId == nil {
		return "", fmt.Errorf("signature missing issuer key id subpacket")
	}

	return formatPGPKeyID(*signature.IssuerKeyId), nil
}

// parseSignature parses an ASCII-armored PGP signature. Expects a single
// packet.
func parseSignature(armoredSignature string) (*packet.Signature, error) {
	block, err := armor.Decode(strings.NewReader(armoredSignature))
	if err != nil {
		return nil, fmt.Errorf("failed to decode armored signature: %w", err)
	}

	p, err := packet.Read(block.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature packet: %w", err)
	}

	s, ok := p.(*packet.Signature)
	if !ok {
		return nil, fmt.Errorf("packet is not a signature")
	}

	// Ensure that the signature was the only packet. Another read should return
	// EOF.
	_, err = packet.Read(block.Body)
	if err != io.EOF {
		return nil, fmt.Errorf("signature contains unexpected packets")
	}

	return s, nil
}

// formatPGPKeyID returns the canonical string representation of a PGP key ID
// (16 hex characters).
func formatPGPKeyID(keyID uint64) string {
	return fmt.Sprintf("%016X", keyID) // note required 0 padding
}

// CheckSignatureByKey checks that `signature` is valid for `payload`
// with the PGP public key in `base64Key`.
func CheckSignatureByKey(base64Key, armoredSignature, payload string) error {
	// Parse the key into a key ring containing only this key.
	keyRing, err := openpgp.ReadKeyRing(base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Key)))
	if err != nil {
		return fmt.Errorf("failed to parse key: %w", err)
	}

	_, err = openpgp.CheckArmoredDetachedSignature(keyRing, strings.NewReader(payload), strings.NewReader(armoredSignature))
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}
