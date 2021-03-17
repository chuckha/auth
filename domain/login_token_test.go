package domain

import (
	"testing"
	"time"
)

func TestNewLoginToken(t *testing.T) {
	_, err := NewLoginToken(&oneTimeToken{}, time.Now().Add(2*time.Second), time.Now())
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotBefore(t *testing.T) {
	token, err := NewLoginToken(&oneTimeToken{}, time.Now().Add(3*time.Second), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if token.notBefore != token.GetNotBefore() {
		t.Fatal("getter is broken on login token")
	}
}
