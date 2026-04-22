package security

import (
	"testing"

	"exam/internal/consts"
)

func TestTokenKey(t *testing.T) {
	if got := tokenKey(consts.UserTypeAdmin, "abc"); got != consts.TokenRedisKeyPrefix+consts.UserTypeTagAdmin+":abc" {
		t.Fatalf("unexpected admin token key: %s", got)
	}
	if got := tokenKey(consts.UserTypeClient, "xyz"); got != consts.TokenRedisKeyPrefix+consts.UserTypeTagClient+":xyz" {
		t.Fatalf("unexpected client token key: %s", got)
	}
}

func TestDecodeTokenPayload(t *testing.T) {
	payload, err := decodeTokenPayload([]byte(`{"user_id":123,"username":"Admin"}`))
	if err != nil {
		t.Fatalf("decode token payload failed: %v", err)
	}
	if payload.UserId != 123 || payload.Username != "Admin" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestDecodeTokenPayloadRejectsInvalidPayload(t *testing.T) {
	if _, err := decodeTokenPayload([]byte(`{"user_id":0,"username":"Admin"}`)); err == nil {
		t.Fatal("expected invalid payload error")
	}
	if _, err := decodeTokenPayload([]byte(`{`)); err == nil {
		t.Fatal("expected malformed payload error")
	}
}
