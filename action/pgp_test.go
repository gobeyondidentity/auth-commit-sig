package action

import (
	"testing"
)

func TestCheckSignatureByKey(t *testing.T) {
	tests := []struct {
		name             string
		base64Key        string
		armoredSignature string
		payload          string
		expectedErr      string
	}{
		{
			name:             "bad_key",
			base64Key:        `nonsense-key`,
			armoredSignature: `nonsense-signature`,
			payload:          `nonsense-payload`,
			expectedErr:      "failed to parse key: illegal base64 data at input byte 0",
		},
		{
			name:             "bad_sig",
			base64Key:        `mFIEYN8gKxMIKoZIzj0DAQcCAwRaSd4BvxkN2rmaRSCd6byYqBxn4jYO0iOsA8RRlC4qZxQ9cD1ssqF5FcriDLLRi0EnB0Jfq1Reo1T2Tn0DUveCtBdUZXN0IDx0ZXN0QGV4YW1wbGUuY29tPohkBBMTCAAWBQJg3yArCRBR9PdBhRpjZQIbAwIZAQAAfJMBAKSHVATg/o2PLo99ueHNkWJ2T7EYe1xGEkrRMB0IoqL3APwMMwNUbW2C715qPBRSRTCObUbQ9XYN0GqKE3LbjNx2uA==`,
			armoredSignature: `nonsense-signature`,
			payload:          `nonsense-payload`,
			expectedErr:      "signature verification failed: EOF",
		},
		{
			name:      "wrong_key",
			base64Key: `mFIEYN8gKxMIKoZIzj0DAQcCAwRaSd4BvxkN2rmaRSCd6byYqBxn4jYO0iOsA8RRlC4qZxQ9cD1ssqF5FcriDLLRi0EnB0Jfq1Reo1T2Tn0DUveCtBdUZXN0IDx0ZXN0QGV4YW1wbGUuY29tPohkBBMTCAAWBQJg3yArCRBR9PdBhRpjZQIbAwIZAQAAfJMBAKSHVATg/o2PLo99ueHNkWJ2T7EYe1xGEkrRMB0IoqL3APwMMwNUbW2C715qPBRSRTCObUbQ9XYN0GqKE3LbjNx2uA==`,
			armoredSignature: `-----BEGIN PGP SIGNATURE-----

iHUEABMIAB0WIQRDF6OV11XR31Mwip/ncqGRwe7exQUCYOroGQAKCRDncqGRwe7e
xQJTAQC1YZqxt3Bf3zkHlUOC9nItItIZF+UZH7B3orT6TEq7yAEAuHKMgVnRXVf5
Qp0Ij4pGLgG+PXYrhj/riYnrRhXwpn4=
=VsS4
-----END PGP SIGNATURE-----
`,
			payload:     "This is a test.\n",
			expectedErr: "signature verification failed: openpgp: signature made by unknown entity",
		},
		{
			name:      "simple",
			base64Key: `mFIEYMI4hhMIKoZIzj0DAQcCAwS+/+HuWZbnBnX6B/lfxZa14RKUvfQKV/gh5Pa0HVRXeBTmNLgsVv7ZDOFf2oLNq0QMYv5B9hK4LSwTogcDuP8atCxGb3JkIEh1cmxleSA8Zm9yZC5odXJsZXlAYmV5b25kaWRlbnRpdHkuY29tPoiQBBMTCAA4FiEEQxejlddV0d9TMIqf53KhkcHu3sUFAmDCOIYCGwMFCwkIBwIGFQoJCAsCBBYCAwECHgECF4AACgkQ53KhkcHu3sVKgQEAp7UdpYP4dSubZQoChsK+QrH96+a4Q3yPyaHsOreDsvIBAP9RcKj9RSVlcDXuUHYfYr25RpvRxAQkuLqbkmusKrdk`,
			armoredSignature: `-----BEGIN PGP SIGNATURE-----

iHUEABMIAB0WIQRDF6OV11XR31Mwip/ncqGRwe7exQUCYOroGQAKCRDncqGRwe7e
xQJTAQC1YZqxt3Bf3zkHlUOC9nItItIZF+UZH7B3orT6TEq7yAEAuHKMgVnRXVf5
Qp0Ij4pGLgG+PXYrhj/riYnrRhXwpn4=
=VsS4
-----END PGP SIGNATURE-----
`,
			payload: "This is a test.\n",
		},
		{
			name:      "good_commit",
			base64Key: `mFIEYMI4hhMIKoZIzj0DAQcCAwS+/+HuWZbnBnX6B/lfxZa14RKUvfQKV/gh5Pa0HVRXeBTmNLgsVv7ZDOFf2oLNq0QMYv5B9hK4LSwTogcDuP8atCxGb3JkIEh1cmxleSA8Zm9yZC5odXJsZXlAYmV5b25kaWRlbnRpdHkuY29tPoiQBBMTCAA4FiEEQxejlddV0d9TMIqf53KhkcHu3sUFAmDCOIYCGwMFCwkIBwIGFQoJCAsCBBYCAwECHgECF4AACgkQ53KhkcHu3sVKgQEAp7UdpYP4dSubZQoChsK+QrH96+a4Q3yPyaHsOreDsvIBAP9RcKj9RSVlcDXuUHYfYr25RpvRxAQkuLqbkmusKrdk`,
			armoredSignature: `-----BEGIN PGP SIGNATURE-----

iJUEABMIAD0WIQRDF6OV11XR31Mwip/ncqGRwe7exQUCYOrmfh8cZm9yZC5odXJs
ZXlAYmV5b25kaWRlbnRpdHkuY29tAAoJEOdyoZHB7t7Fs/MA/jfzo9cigGqEvmVz
YIKMWCp0G4FD2Gp54QLlr5osrtf9AP98F4zFigLfCQG7ria/lxvyjHx0khmpFpO4
Q/RgcSQ+iw==
=gFHw
-----END PGP SIGNATURE-----
`,
			payload: `tree 3657f03df2a230f4c0f8682e515efd5bdc030cb9
parent cf52d82d9ea21d5c9d174b21936aed4d01fbbd25
author Ford Hurley <ford.hurley@beyondidentity.com> 1626006821 -0400
committer Ford Hurley <ford.hurley@beyondidentity.com> 1626007165 -0400

Add usage to README
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSignatureByKey(tt.base64Key, tt.armoredSignature, tt.payload)
			assertEqualErr(t, tt.expectedErr, err)
		})
	}
}

func TestParseSignatureIssuerKeyID(t *testing.T) {
	tests := []struct {
		name             string
		armoredSignature string
		expected         string
		expectedErr      string
	}{
		{
			name: "simple",
			armoredSignature: `-----BEGIN PGP SIGNATURE-----

iJUEABMIAD0WIQRDF6OV11XR31Mwip/ncqGRwe7exQUCYOrmfh8cZm9yZC5odXJs
ZXlAYmV5b25kaWRlbnRpdHkuY29tAAoJEOdyoZHB7t7Fs/MA/jfzo9cigGqEvmVz
YIKMWCp0G4FD2Gp54QLlr5osrtf9AP98F4zFigLfCQG7ria/lxvyjHx0khmpFpO4
Q/RgcSQ+iw==
=gFHw
-----END PGP SIGNATURE-----
`,
			expected: "E772A191C1EEDEC5",
		},
		{
			name:             "nonsense",
			armoredSignature: `nonsense`,
			expectedErr:      "failed to parse signature: failed to decode armored signature: EOF",
		},
		{
			name: "armored_nonsense",
			armoredSignature: `-----BEGIN PGP SIGNATURE-----

nonsense
-----END PGP SIGNATURE-----
`,
			expectedErr: "failed to parse signature: failed to read signature packet: unexpected EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSignatureIssuerKeyID(tt.armoredSignature)
			assertEqualErr(t, tt.expectedErr, err)

			if err != nil && got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func assertEqualErr(t *testing.T, expected string, err error) {
	t.Helper()

	if expected == "" {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		return
	}

	if err == nil {
		t.Errorf("expected error %v, got nil", expected)
		return
	}

	if err.Error() != expected {
		t.Errorf("expected error %v, got %v", expected, err)
	}
}
