package log_message

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

func (r *DatabaseRepository) FindByID(ctx context.Context, ID *int) (*Model, error) {
	var payout Model
	query := r.db.NewSelect().
		Model(&payout).
		Where("id = ?", ID)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, terrors.New(terrors.ErrNotFound, "Payout information not found", map[string]string{})
		}
		return nil, fmt.Errorf("payout_repository: Error while searching for countrysvc -> %v", err)
	}

	return &payout, nil

}

// Create Handles the creation of a new payout record on a database
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

func (r *DatabaseRepository) Update(ctx context.Context, id *int, model *Model) error {
	query := r.db.NewUpdate().
		Table("log_messages").
		Set("message = ?", model.Message).
		Where("id = ?", id)

	// Execute the query
	res, err := query.Exec(ctx)
	if err != nil {
		return err
	}

	// Retrieve total of affected rows
	updatedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("log_message_repo: error while retrieving rows affected: %w", err)
	}

	// Validate that one record was updated
	if updatedRows != 1 {
		return terrors.New(terrors.ErrNotFound, "Record not found", map[string]string{})
	}

	return nil
}
