package main

import (
	"fmt"
	"net/http"

	"github.com/chuckha/auth/app"
	"github.com/chuckha/auth/infrastructure/comms/messages"
	"github.com/chuckha/auth/infrastructure/comms/terminal"
	"github.com/chuckha/auth/infrastructure/paseto"
	"github.com/chuckha/auth/infrastructure/secret"
	"github.com/chuckha/auth/infrastructure/sqlite"
	"github.com/chuckha/auth/infrastructure/uuid"
	"github.com/chuckha/auth/usecases"
)

func main() {
	store, err := sqlite.NewSQLiteStore("auth-db")
	if err != nil {
		panic(err)
	}
	if err := store.Initialize("cmd/auth/migrations"); err != nil {
		panic(err)
	}

	keyGetter := &secret.KeyGetter{}
	encdec := paseto.NewPASETOEncDec()
	messageGen := &messages.Generator{}
	messageSender := &terminal.Client{}

	controller := app.NewController(&app.AuthUseCases{
		DoLogin: &usecases.DoLogin{
			IDGenerator:       &uuid.UUID{},
			SecretKeyGetter:   keyGetter,
			Decoder:           encdec,
			TokenRepository:   store,
			SessionRepository: store,
		},
		LoginMessageSender: &usecases.LoginMessageSender{
			IDGenerator:      &uuid.UUID{},
			SecretKeyGetter:  keyGetter,
			Encoder:          encdec,
			MessageGenerator: messageGen,
			MessageSender:    messageSender,
			TokenRepository:  store,
			UserRepository:   store,
		},
	}, &app.UseCasesAdapter{}, &app.UseCasesPresenter{})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Printf("error parsing form: %+v\n", err)
			return
		}
		email := r.Form.Get("email")
		_, err := controller.SendLoginMessage(&app.SendLoginMessageInput{LoginMessageDestination: email})
		if err != nil {
			fmt.Printf("error sending login message: %+v\n", err)
			return
		}
	})

	http.HandleFunc("/validate_link", func(w http.ResponseWriter, r *http.Request) {
		tkn := r.URL.Query().Get("token")
		out, err := controller.Login(&app.LoginInput{EncodedToken: tkn})
		if err != nil {
			fmt.Printf("error sending login message: %+v\n", err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: out.SessionID,
			Path:  "/",
		})
	})
	fmt.Println("listening on 8888")
	panic(http.ListenAndServe(":8888", nil))
}
