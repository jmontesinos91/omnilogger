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

// FindByID finds a log message by its ID
func (r *DatabaseRepository) FindByID(ctx context.Context, ID *int) (*Model, error) {
	var logMessage Model
	query := r.db.NewSelect().
		Model(&logMessage).
		Where("id = ?", ID)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, terrors.New(terrors.ErrNotFound, "Log Message information not found", map[string]string{})
		}
		return nil, fmt.Errorf("log_message_repository: Error while searching for logmessage -> %v", err)
	}

	return &logMessage, nil

}

// FindByIDAndLang finds a log message by its ID and language
func (r *DatabaseRepository) FindByIDAndLang(ctx context.Context, ID *int, lang string) (*Model, error) {
	var logMessage Model
	query := r.db.NewSelect().
		Model(&logMessage).
		Where("id = ?", ID).
		Where("lang = ?", lang)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, terrors.New(terrors.ErrNotFound, "Log Message information not found", map[string]string{})
		}
		return nil, fmt.Errorf("log_message_repository: Error while searching for logmessage -> %v", err)
	}

	return &logMessage, nil
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

// Update Handles the update of an existing log message record on a database
func (r *DatabaseRepository) Update(ctx context.Context, id *int, lang string, model *Model) error {
	query := r.db.NewUpdate().
		Table("log_messages").
		Set("message = ?", model.Message).
		Set("lang = ?", model.Lang).
		Where("id = ?", id).
		Where("lang = ?", lang)

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

// Retrieve retrieves log messages based on the provided filter
func (r *DatabaseRepository) Retrieve(ctx context.Context, filter Filter) ([]Model, int, error) {
	var model []Model
	query := r.db.NewSelect().Model(&model).
		Order("id ASC").
		Limit(filter.Size).
		Offset(filter.From - 1)

	if filter.ID != nil {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Lang != "" {
		query = query.Where("lang = ?", filter.Lang)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, 0, err
	}

	count, _ := query.Count(ctx)

	return model, count, nil
}

// DeleteLang deletes a log message by its ID and language
func (r *DatabaseRepository) DeleteLang(ctx context.Context, id *int, lang string) error {
	var logMessage Model
	query := r.db.NewDelete().
		Model(&logMessage).
		Where("id = ?", id).
		Where("lang = ?", lang)

	if _, err := query.Exec(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return terrors.New(terrors.ErrNotFound, "Log message information not found", map[string]string{})
		}
		return fmt.Errorf("log_message_repository: Error while deleting for log_message_svc -> %w", err)
	}

	return nil
}
