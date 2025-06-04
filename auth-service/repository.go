package authservice

import (
	"context"
	"database/sql"
)

type Repository interface{}

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

func (r *userDatabase) NewCard(ctx context.Context, card CardMetadata) error {
	q := `
		INSERT INTO cards(
			user_id,
			card_number,
			brand,
			expiry_month,
			expiry_year,
		) VALUES ($1,$2,$3,$4,$5)
	`
	_, err := r.db.ExecContext(ctx, q, card.UID, card.Number, card.Brand, card.ExpiryMonth, card.ExpiryYear)
	return err
}

func (r *userDatabase) GetUser(ctx context.Context, id string) (*UserResponseModel, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var txErr error // TRANSECTION ERROR
	defer func() {
		if txErr != nil {
			tx.Rollback()
		}
	}()

	var q string // QUERY STRING

	user := new(UserModel) // USER FIELD
	user.ID = id
	q = `
		SELECT username, password, email, email,r_token
		FROM users 
		WHERE id = $1 
	`
	txErr = tx.QueryRowContext(ctx, q, user.ID).Scan(
		&user.Username,
		&user.Password,
		&user.Email,
		&user.RefreshToken,
	)
	if txErr != nil {
		return nil, txErr
	}

	cards := new([]CardsResponseMetadata) // CARD FIELD
	q = `
		SELECT card_number, brand, expiry_month, expiry_year
		FROM cards 
		WHERE user_id = $1 
	`
	cardRows, txErr := tx.QueryContext(ctx, q, id)
	if txErr != nil {
		return nil, txErr
	}
	defer cardRows.Close()
	for cardRows.Next() {
		c := CardsResponseMetadata{}
		txErr = cardRows.Scan(
			&c.Number,
			&c.Brand,
			&c.ExpiryMonth,
			&c.ExpiryYear,
		)
		if txErr != nil {
			return nil, txErr
		}
		*cards = append(*cards, c)
	}

	if txErr = tx.Commit(); txErr != nil {
		return nil, txErr
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

func (r *userDatabase) UpdateCard(ctx context.Context, card CardMetadata) error {
	q := `
		UPDATE cards
		SET
			card_number = $2,
			brand = $3,
			expiry_month = $4,
			expiry_year = $5
		WHERE user_id = $1	
	`
	_, err := r.db.ExecContext(ctx, q, card.UID, card.Number, card.Brand, card.ExpiryMonth, card.ExpiryYear)
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

func (r *userDatabase) DeleteCard(ctx context.Context, uid string) error {
	q := `
		DELETE FROM cards
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, q, uid)
	return err
}
