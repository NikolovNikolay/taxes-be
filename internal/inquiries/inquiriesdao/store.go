package inquiriesdao

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"taxes-be/internal/core"
	"taxes-be/internal/daoutil"
	"taxes-be/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) FindInquiry(ctx context.Context, id uuid.UUID) (*models.Inquiry, error) {
	var inquiry *models.Inquiry
	err := daoutil.EnsureTransaction(ctx, s.db, func(tx *sql.Tx) error {
		inquiryModel, err := models.Inquiries(
			models.InquiryWhere.ID.EQ(id.String()),
		).One(ctx, tx)
		if err != nil {
			if err == sql.ErrNoRows {
				return core.ErrNotFound(err)
			}

			return err
		}
		inquiry = inquiryModel

		return nil
	})
	if err != nil {
		return nil, err
	}

	return inquiry, nil
}

func (s *Store) UpsertInquiry(ctx context.Context, inquiry *models.Inquiry) error {
	err := daoutil.EnsureTransaction(ctx, s.db, func(tx *sql.Tx) error {
		err := inquiry.Upsert(ctx, tx, true, []string{
			models.CouponColumns.ID,
		}, boil.Infer(), boil.Infer())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) InTransaction(ctx context.Context, retryableFn func(ctx context.Context) error) error {
	return daoutil.InTransaction(ctx, s.db, retryableFn)
}
