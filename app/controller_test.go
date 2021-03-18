package app

import (
	"testing"
	"time"

	"github.com/chuckha/auth/domain"
	"github.com/chuckha/auth/usecases"
)

type fakesSeshRep struct {
	out *domain.Session
	err error
}

func (f *fakesSeshRep) LookupSession(_ string) (*domain.Session, error) {
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
					out: &domain.Session{UID: "abc", ID: "def", Expires: time.Now().Add(1 * time.Second)},
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
