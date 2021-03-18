package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/chuckha/auth/domain"
)

const (
	SessionsTableName = "sessions"
	SessionFields     = "id, user_id, expires"

	TokensTableName = "tokens"
	TokenFields     = "token, user_id, expires"

	UsersTableName = "users"
	UserFields     = "id"
)

type Store struct {
	DB *sql.DB
}

func NewSQLiteStore(database string) (*Store, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Store{db}, nil
}

func (s *Store) LookupSession(id string) (*domain.Session, error) {
	out := &domain.Session{}
	expires := ""
	err := s.DB.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE id = ?`, SessionFields, SessionsTableName), id).
		Scan(&out.ID, &out.UID, &expires)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	t, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	out.Expires = t
	return out, nil
}

func (s *Store) GetSession(uid, id string) (*domain.Session, error) {
	out := &domain.Session{}
	expires := ""
	err := s.DB.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE id = ? AND user_id = ?`, SessionFields, SessionsTableName), id, uid).
		Scan(&out.ID, &out.UID, &expires)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	t, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	out.Expires = t
	return out, nil
}

func (s *Store) SaveSession(session *domain.Session) error {
	_, err := s.DB.Exec(fmt.Sprintf(`INSERT INTO %s (%s) VALUES (?, ?, ?)`, SessionsTableName, SessionFields),
		session.ID, session.UID, session.Expires.Format(time.RFC3339))
	return errors.WithStack(err)
}

func (s *Store) GetToken(uid, token string) (*domain.OneTimeToken, error) {
	out := &domain.OneTimeToken{}
	expires := ""
	err := s.DB.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE token = ? AND user_id = ?`, TokenFields, TokensTableName), token, uid).
		Scan(&out.Token, &out.UserID, &expires)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	t, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	out.Expires = t
	return out, nil
}

func (s *Store) SaveToken(token *domain.OneTimeToken) error {
	_, err := s.DB.Exec(fmt.Sprintf(`INSERT INTO %s (%s) VALUES (?, ?, ?)`, TokensTableName, TokenFields), token.Token, token.UserID, token.Expires.Format(time.RFC3339))
	return errors.WithStack(err)
}

func (s *Store) DeleteToken(uid, token string) error {
	_, err := s.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE token = ? AND user_id = ?`, TokensTableName), token, uid)
	return errors.WithStack(err)
}

func (s *Store) GetUser(uid string) (*domain.User, error) {
	out := &domain.User{}
	err := s.DB.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE id = ?`, UserFields, UsersTableName), uid).Scan(&out.ID)
	return out, errors.WithStack(err)
}

func (s *Store) CreateUser(user *domain.User) error {
	_, err := s.DB.Exec(fmt.Sprintf(`INSERT INTO %s (%s) VALUES (?)`, UsersTableName, UserFields), user.ID)
	return errors.WithStack(err)
}
