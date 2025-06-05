package authservice

import (
	"context"
	"database/sql"
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

func NewUserDatabase(dsn string) (Repository, error) {
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
	q := `
		INSERT INTO users(
			id,
			username,
			password,
			email,
			r_token,
		) VALUES ($1,$2,$3,$4,$5)
	`
	_, err := r.db.ExecContext(ctx, q, user.ID, user.Username, user.Password, user.Email, user.RefreshToken)
	return err
}

func (r *userDatabase) GetUser(ctx context.Context, id string) (*UserResponseModel, error) {
	var q string           // QUERY STRING
	user := new(UserModel) // USER FIELD
	user.ID = id
	q = `
		SELECT username, password, email, email,r_token
		FROM users 
		WHERE id = $1 
	`
	err := r.db.QueryRowContext(ctx, q, user.ID).Scan(
		&user.Username,
		&user.Password,
		&user.Email,
		&user.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	cards, err := r.GetCardsById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserResponseModel{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Token: TokenMetadata{
			RefreshToken: user.RefreshToken,
		},
		Cards: *cards,
	}, nil
}

func (r *userDatabase) GetUserByEmail(ctx context.Context, email string) (*UserResponseModel, error) {
	user := new(UserModel) // USER FIELD
	q := `
		SELECT username, password, email, email,r_token
		FROM users 
		WHERE email = $1 
	`
	err := r.db.QueryRowContext(ctx, q, user.Email).Scan(
		&user.Username,
		&user.Password,
		&user.Email,
		&user.RefreshToken,
	)
	if err != nil {
		return nil, err
	}
	return r.GetUser(ctx, user.ID)
}

func (r *userDatabase) UpdateUser(ctx context.Context, user UserModel) error {
	q := `
		UPDATE users
		SET
			username = $2,
			password = $3,
			email = $4,
			r_token = $5,
		WHERE id = $1	
	`
	_, err := r.db.ExecContext(ctx, q, user.ID, user.Username, user.Password, user.Email, user.RefreshToken)
	return err
}

func (r *userDatabase) UpdateToken(ctx context.Context, id string, token string) error {
	q := `
		UPDATE users
		SET r_token = $2,
		WHERE id = $1	
	`
	_, err := r.db.ExecContext(ctx, q, id, token)
	return err
}

func (r *userDatabase) DeleteUser(ctx context.Context, id string) error {
	q := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

// CARD SECTION
func (r *userDatabase) NewCard(ctx context.Context, card CardMetadata) (*[]CardsResponseMetadata, error) {
	q := `
		INSERT INTO cards(
			user_id,
			card_number,
			brand,
			expiry_month,
			expiry_year,
		) VALUES ($1,$2,$3,$4,$5)
	`
	if _, err := r.db.ExecContext(ctx, q, card.UID, card.Number, card.Brand, card.ExpiryMonth, card.ExpiryYear); err != nil {
		return nil, err
	}

	return r.GetCardsById(ctx, card.UID)
}

func (r *userDatabase) GetCardsById(ctx context.Context, id string) (*[]CardsResponseMetadata, error) {
	cards := new([]CardsResponseMetadata) // CARD FIELD
	q := `
		SELECT id, card_number, brand, expiry_month, expiry_year
		FROM cards 
		WHERE user_id = $1 
	`
	cardRows, txErr := r.db.QueryContext(ctx, q, id)
	if txErr != nil {
		return nil, txErr
	}
	defer cardRows.Close()
	for cardRows.Next() {
		c := CardsResponseMetadata{}
		txErr = cardRows.Scan(
			&c.ID,
			&c.Number,
			&c.Brand,
			&c.ExpiryMonth,
			&c.ExpiryYear,
		)
		if txErr != nil {
			return nil, txErr
		}
		c.Number = hideNumber(c.Number)
		*cards = append(*cards, c)
	}

	return cards, nil
}

func (r *userDatabase) DeleteAllCard(ctx context.Context, uid string) error {
	q := `
		DELETE FROM cards
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, q, uid)
	return err
}

func (r *userDatabase) DeleteCard(ctx context.Context, uid string, id string) error { // ID -> CARD_ID .. UID -> USER_ID
	q := `
		DELETE FROM cards
		WHERE user_id = $1 AND id = $2
	`
	_, err := r.db.ExecContext(ctx, q, uid, id)
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
