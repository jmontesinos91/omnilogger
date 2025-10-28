package logs

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/terrors"
	"github.com/uptrace/bun"
)

// DatabaseRepository struct
type DatabaseRepository struct {
	log *logger.ContextLogger
	db  *bun.DB
}

// NewDatabaseRepository creates an instance of DatabaseRepository
func NewDatabaseRepository(l *logger.ContextLogger, conn *bun.DB) *DatabaseRepository {
	return &DatabaseRepository{
		log: l,
		db:  conn,
	}
}

func (r *DatabaseRepository) FindByID(ctx context.Context, ID *string) (*Model, error) {
	var payout Model
	query := r.db.NewSelect().
		Model(&payout).
		Where("id = ?", ID)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, terrors.New(terrors.ErrNotFound, "Log information not found", map[string]string{})
		}
		return nil, fmt.Errorf("payout_repository: Error while searching for countrysvc -> %v", err)
	}

	return &payout, nil
}

// Create Handles the creation of a new log record on database
func (r *DatabaseRepository) Create(ctx context.Context, model *Model) error {
	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	// Handling error
	if err != nil {
		return err
	}
	return nil
}
