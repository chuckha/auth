package app

import (
	"testing"
	"time"

	"github.com/chuckha/auth/usecases"
	"github.com/chuckha/auth/usecases/dto"
)

type fakesSeshRep struct {
	out *dto.Session
	err error
}

func (f *fakesSeshRep) LookupSession(_ string) (*dto.Session, error) {
	return f.out, f.err
}

func TestController_CheckValidSession(t *testing.T) {
	c := &Controller{
		Adapter: &UseCasesAdapter{},
		UseCases: &AuthUseCases{
			DoLogin:            nil,
			LoginMessageSender: nil,
			SessionChecker: &usecases.SessionChecker{
				LookupSessionRepository: &fakesSeshRep{
					out: &dto.Session{UserID: "abc", ID: "def", Expires: time.Now().Add(1 * time.Second)},
					err: nil,
				},
			},
		},
		Presenter: &UseCasesPresenter{},
	}
	out, err := c.CheckValidSession(&CheckValidSessionInput{SessionID: "def"})
	if err != nil {
		t.Fatal(err)
	}
	if out.UID != "abc" {
		t.Fatalf("uid needs to be %q not %q", "abc", out.UID)
	}
}
