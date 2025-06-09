package authservice

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	LOCATION = "DB_REPOSITORY"
)

type Repository interface {
	NewUser(ctx context.Context, user UserModel) error
	NewCard(ctx context.Context, card CardMetadata) (*[]CardsResponseMetadata, error)
	GetUser(ctx context.Context, id string) (*UserResponseModel, error)
	GetCardsById(ctx context.Context, id string) (*[]CardsResponseMetadata, error)
	GetUserByEmail(ctx context.Context, email string) (*UserResponseModel, error)
	UpdateUser(ctx context.Context, user UserModel) error
	UpdateToken(ctx context.Context, id string, token string) error
	DeleteUser(ctx context.Context, id string) error
	DeleteAllCard(ctx context.Context, uid string) error
	DeleteCard(ctx context.Context, uid string, id string) error
}

type userDatabase struct {
	db *sql.DB
}

func NewRepository(dsn string) (Repository, error) {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &userDatabase{db: db}, nil
}

func (r *userDatabase) NewUser(ctx context.Context, user UserModel) error {
	var err error
	defer func() {
		log.Println(LOCATION, user)
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	q := `
		INSERT INTO users(
			id,
			username,
			password,
			email,
			r_token
		) VALUES ($1,$2,$3,$4,$5)
	`
	_, err = r.db.ExecContext(ctx, q, user.ID, user.Username, user.Password, user.Email, user.RefreshToken)
	return err
}

func (r *userDatabase) GetUser(ctx context.Context, id string) (*UserResponseModel, error) {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	var q string           // QUERY STRING
	user := new(UserModel) // USER FIELD
	user.ID = id
	q = `
		SELECT username, password, email,r_token
		FROM users 
		WHERE id = $1 
	`
	err = r.db.QueryRowContext(ctx, q, id).Scan(
		&user.Username,
		&user.Password,
		&user.Email,
		&user.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	// cards, err := r.GetCardsById(ctx, id)
	// if err != nil {
	// 	return nil, err
	// }
	return &UserResponseModel{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Token: TokenMetadata{
			RefreshToken: user.RefreshToken,
		},
		Cards: nil,
	}, nil
}

func (r *userDatabase) GetUserByEmail(ctx context.Context, email string) (*UserResponseModel, error) {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	user := new(UserModel) // USER FIELD
	q := `
		SELECT id
		FROM users 
		WHERE email = $1
 	`
	err = r.db.QueryRowContext(ctx, q, email).Scan(
		&user.ID,
	)
	if err != nil {
		return nil, err
	}
	return r.GetUser(ctx, user.ID)
}

func (r *userDatabase) UpdateUser(ctx context.Context, user UserModel) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	q := `
		UPDATE users
		SET
			username = $2,
			password = $3,
			email = $4,
			r_token = $5
		WHERE id = $1	
	`
	_, err = r.db.ExecContext(ctx, q, user.ID, user.Username, user.Password, user.Email, user.RefreshToken)
	return err
}

func (r *userDatabase) UpdateToken(ctx context.Context, id string, token string) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	q := `
		UPDATE users
		SET r_token = $2
		WHERE id = $1	
	`
	_, err = r.db.ExecContext(ctx, q, id, token)
	return err
}

func (r *userDatabase) DeleteUser(ctx context.Context, id string) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	q := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err = r.db.ExecContext(ctx, q, id)
	return err
}

// CARD SECTION
func (r *userDatabase) NewCard(ctx context.Context, card CardMetadata) (*[]CardsResponseMetadata, error) {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	q := `
		INSERT INTO cards(
			user_id,
			card_number,
			brand,
			expiry_month,
			expiry_year
		) VALUES ($1,$2,$3,$4,$5)
	`
	if _, err = r.db.ExecContext(ctx, q, card.UID, card.Number, card.Brand, card.ExpiryMonth, card.ExpiryYear); err != nil {
		log.Println(err)
		return nil, err
	}

	return r.GetCardsById(ctx, card.UID)
}

func (r *userDatabase) GetCardsById(ctx context.Context, uid string) (*[]CardsResponseMetadata, error) {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	cards := new([]CardsResponseMetadata) // CARD FIELD
	q := `
		SELECT id, card_number, brand, expiry_month, expiry_year
		FROM cards 
		WHERE user_id = $1 
	`
	cardRows, err := r.db.QueryContext(ctx, q, uid)
	if err != nil {
		return nil, err
	}
	defer cardRows.Close()
	for cardRows.Next() {
		c := CardsResponseMetadata{}
		err = cardRows.Scan(
			&c.ID,
			&c.Number,
			&c.Brand,
			&c.ExpiryMonth,
			&c.ExpiryYear,
		)
		if err != nil {
			return nil, err
		}
		c.Number = hideNumber(c.Number)
		*cards = append(*cards, c)
	}

	return cards, nil
}

func (r *userDatabase) DeleteAllCard(ctx context.Context, uid string) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	q := `
		DELETE FROM cards
		WHERE user_id = $1
	`
	_, err = r.db.ExecContext(ctx, q, uid)
	return err
}

func (r *userDatabase) DeleteCard(ctx context.Context, uid string, id string) error { // ID -> CARD_ID .. UID -> USER_ID
	var err error
	defer func() {
		if err != nil {
			log.Println(LOCATION, err)
		}
	}()
	q := `
		DELETE FROM cards
		WHERE user_id = $1 AND id = $2
	`
	_, err = r.db.ExecContext(ctx, q, uid, id)
	return err
}

func hideNumber(number string) string {
	size := len(number) - 4
	pref := number[:4]
	suf := ""
	for i := 0; i < size; i++ {
		suf += " *"
	}
	return pref + suf
}
