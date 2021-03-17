package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/chuckha/migrations"
	"github.com/pkg/errors"

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
	db, err := sql.Open("sqlite3", "auth")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := Initialize(db); err != nil {
		panic(err)
	}

	keyGetter := &secret.KeyGetter{}
	store := &sqlite.SQLite{db}
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

// Initialize manages migrations for the test db
func Initialize(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (number int)`)
	if err != nil {
		return errors.WithStack(err)
	}
	row := db.QueryRow(`SELECT number FROM migrations ORDER BY number DESC LIMIT 1`)
	if err := row.Err(); err != nil {
		return errors.WithStack(err)
	}
	latest := 0

	if err := row.Scan(&latest); err != nil {
		if err != sql.ErrNoRows {
			return errors.WithStack(err)
		}
	}
	// printlnf := func(a string, args ...interface{}) {
	// 	fmt.Printf(a+"\n", args...)
	// }
	migs := migrations.FromDir("cmd/auth/migrations")
	for _, migration := range migs[:latest] {
		fmt.Println("Already ran migration: %d", migration.Order)
	}
	for _, migration := range migs[latest:] {
		fmt.Printf("Running migration %d\n", migration.Order)
		if _, err := db.Exec(migration.Up); err != nil {
			return err
		}
		if _, err := db.Exec(`INSERT INTO migrations (number) VALUES (?)`, migration.Order); err != nil {
			return err
		}
	}
	return nil
}
