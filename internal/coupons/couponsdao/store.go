package couponsdao

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

func (s *Store) FindCoupon(ctx context.Context, id uuid.UUID) (*models.Coupon, error) {
	var coupon *models.Coupon
	err := daoutil.EnsureTransaction(ctx, s.db, func(tx *sql.Tx) error {
		couponModel, err := models.FindCoupon(ctx, tx, id.String())
		if err != nil {
			if err == sql.ErrNoRows {
				return core.ErrNotFound(err)
			}

			return err
		}
		coupon = couponModel

		return nil
	})
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (s *Store) FindCouponByRequestID(ctx context.Context, id uuid.UUID) (*models.Coupon, error) {
	var coupon *models.Coupon
	err := daoutil.EnsureTransaction(ctx, s.db, func(tx *sql.Tx) error {
		couponModel, err := models.Coupons(
			models.CouponWhere.ParentRequestID.EQ(id.String()),
		).One(ctx, tx)

		if err != nil {
			if err == sql.ErrNoRows {
				return core.ErrNotFound(err)
			}

			return err
		}
		coupon = couponModel

		return nil
	})
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (s *Store) UpsertCoupon(ctx context.Context, coupon *models.Coupon) error {
	err := daoutil.EnsureTransaction(ctx, s.db, func(tx *sql.Tx) error {
		err := coupon.Upsert(ctx, tx, true, []string{
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
