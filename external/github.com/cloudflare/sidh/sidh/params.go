package sidh

import (
	. "v2ray.com/core/external/github.com/cloudflare/sidh/internal/isogeny"
	p503 "v2ray.com/core/external/github.com/cloudflare/sidh/p503"
	p751 "v2ray.com/core/external/github.com/cloudflare/sidh/p751"
)

// Keeps mapping: SIDH prime field ID to domain parameters
var sidhParams = make(map[uint8]SidhParams)

// Params returns domain parameters corresponding to finite field and identified by
// `id` provieded by the caller. Function panics in case `id` wasn't registered earlier.
func Params(id uint8) *SidhParams {
	if val, ok := sidhParams[id]; ok {
		return &val
	}
	panic("sidh: SIDH Params ID unregistered")
}

func init() {
	p503 := SidhParams{
		Id:               FP_503,
		PublicKeySize:    p503.P503_PublicKeySize,
		SharedSecretSize: p503.P503_SharedSecretSize,
		A: DomainParams{
			Affine_P:        p503.P503_affine_PA,
			Affine_Q:        p503.P503_affine_QA,
			Affine_R:        p503.P503_affine_RA,
			SecretBitLen:    p503.P503_SecretBitLenA,
			SecretByteLen:   uint((p503.P503_SecretBitLenA + 7) / 8),
			IsogenyStrategy: p503.P503_AliceIsogenyStrategy[:],
		},
		B: DomainParams{
			Affine_P:        p503.P503_affine_PB,
			Affine_Q:        p503.P503_affine_QB,
			Affine_R:        p503.P503_affine_RB,
			SecretBitLen:    p503.P503_SecretBitLenB,
			SecretByteLen:   uint((p503.P503_SecretBitLenB + 7) / 8),
			IsogenyStrategy: p503.P503_BobIsogenyStrategy[:],
		},
		OneFp2:  p503.P503_OneFp2,
		HalfFp2: p503.P503_HalfFp2,
		MsgLen:  24,
		// SIKEp751 provides 128 bit of classical security ([SIKE], 5.1)
		KemSize: 16,
		Bytelen: p503.P503_Bytelen,
		Op:      p503.FieldOperations(),
	}

	p751 := SidhParams{
		Id:               FP_751,
		PublicKeySize:    p751.P751_PublicKeySize,
		SharedSecretSize: p751.P751_SharedSecretSize,
		A: DomainParams{
			Affine_P:        p751.P751_affine_PA,
			Affine_Q:        p751.P751_affine_QA,
			Affine_R:        p751.P751_affine_RA,
			IsogenyStrategy: p751.P751_AliceIsogenyStrategy[:],
			SecretBitLen:    p751.P751_SecretBitLenA,
			SecretByteLen:   uint((p751.P751_SecretBitLenA + 7) / 8),
		},
		B: DomainParams{
			Affine_P:        p751.P751_affine_PB,
			Affine_Q:        p751.P751_affine_QB,
			Affine_R:        p751.P751_affine_RB,
			IsogenyStrategy: p751.P751_BobIsogenyStrategy[:],
			SecretBitLen:    p751.P751_SecretBitLenB,
			SecretByteLen:   uint((p751.P751_SecretBitLenB + 7) / 8),
		},
		OneFp2:  p751.P751_OneFp2,
		HalfFp2: p751.P751_HalfFp2,
		MsgLen:  32,
		// SIKEp751 provides 192 bit of classical security ([SIKE], 5.1)
		KemSize: 24,
		Bytelen: p751.P751_Bytelen,
		Op:      p751.FieldOperations(),
	}

	sidhParams[FP_503] = p503
	sidhParams[FP_751] = p751
}
