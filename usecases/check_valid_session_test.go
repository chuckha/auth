package usecases

import (
	"testing"
	"time"

	"github.com/chuckha/auth/domain"
)

type fakesSeshRep struct {
	out *domain.Session
	err error
}

func (f *fakesSeshRep) LookupSession(id string) (*domain.Session, error) {
	return f.out, f.err
}

func TestCheckValidSession_CheckValidSession(t *testing.T) {
	cvs := &SessionChecker{
		LookupSessionRepository: &fakesSeshRep{
			out: &domain.Session{
				ID:      "def",
				UID:     "abc",
				Expires: time.Now().Add(10 * time.Minute),
			},
			err: nil,
		},
	}
	out, err := cvs.CheckValidSession(&ValidSessionInput{
		SID: "def",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.UID != "abc" {
		t.Fatal("uid doesn't match")
	}
}
