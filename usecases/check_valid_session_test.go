package usecases

import (
	"testing"
	"time"

	"github.com/chuckha/auth/usecases/dto"
)

type fakesSeshRep struct {
	out *dto.Session
	err error
}

func (f *fakesSeshRep) GetSession(uid, id string) (*dto.Session, error) {
	return f.out, f.err
}

func TestCheckValidSession_CheckValidSession(t *testing.T) {
	cvs := &SessionChecker{
		GetSessionRepository: &fakesSeshRep{
			out: &dto.Session{
				ID:      "def",
				UserID:  "abc",
				Expires: time.Now().Add(10 * time.Minute),
			},
			err: nil,
		},
	}
	out, err := cvs.CheckValidSession(&ValidSessionInput{
		UID: "abc",
		SID: "def",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.UID != "abc" {
		t.Fatal("uid doesn't match")
	}
	if out.SID != "def" {
		t.Fatal("sid doesn't match")
	}
}
