package logic

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func sign(t *testing.T, secret, canonical string) string {
	t.Helper()
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(canonical))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifySignAcceptsValidSignature(t *testing.T) {
	const secret = "test-secret"
	s := sign(t, secret, "payOrderNo=990010")
	if !verifySign(secret, 990010, s) {
		t.Fatal("valid signature should pass")
	}
}

func TestVerifySignRejectsTampered(t *testing.T) {
	const secret = "test-secret"
	s := sign(t, secret, "payOrderNo=990010")
	if verifySign(secret, 990011, s) {
		t.Fatal("signature for another order should fail")
	}
	if verifySign("other-secret", 990010, s) {
		t.Fatal("signature with wrong secret should fail")
	}
	if verifySign(secret, 990010, s[:len(s)-2]+"00") {
		t.Fatal("tampered signature should fail")
	}
}

func TestVerifySignRejectsEmpty(t *testing.T) {
	if verifySign("", 1, "x") || verifySign("s", 1, "") {
		t.Fatal("empty secret or sign should fail")
	}
}

// 锁定签名规则文档中的 canonical 形态：payOrderNo=<no>
func TestSignCanonicalForm(t *testing.T) {
	const secret = "k"
	s := sign(t, secret, "payOrderNo=42")
	if !verifySign(secret, 42, s) {
		t.Fatal("canonical form payOrderNo=<no> should be the signing input")
	}
	if verifySign(secret, 42, sign(t, secret, "payOrderNo=42&result=SUCCESS")) {
		t.Fatal("extra params must NOT be part of canonical form")
	}
}
