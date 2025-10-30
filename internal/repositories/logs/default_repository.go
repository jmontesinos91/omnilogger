package logs

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/sts"
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

func (r *DatabaseRepository) FindByID(ctx context.Context, ID *string, filter Filter) (*Model, error) {
	var payout Model
	query := r.db.NewSelect().
		Model(&payout).
		Relation("LogMessage", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("model.lang = ?", filter.Lang)
		}).
		Where("id = ?", ID)

	if err := query.Scan(ctx); err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, terrors.New(terrors.ErrNotFound, "Log information not found", map[string]string{})
		}
		return nil, fmt.Errorf("payout_repository: Error while searching for countrysvc -> %v", err)
	}

	return &payout, nil
}

// Create Handles the creation of a new log record on a database
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

func (r *DatabaseRepository) Retrieve(ctx context.Context, filter Filter) ([]Model, int, error) {
	claims := ctx.Value(&sts.Claim).(sts.Claims)
	userTenantsID := claims.Tenants

	var model []Model
	query := r.db.NewSelect().Model(&model).
		Relation("LogMessage", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("model.lang = ?", filter.Lang)
		}).
		Order("created_at DESC").
		Limit(filter.Size).
		Offset(filter.From - 1)

	if len(filter.Message) > 0 {
		query = query.Where("message in (?)", bun.In(filter.Message))
	}

	if len(filter.Level) > 0 {
		query = query.Where("level in (?)", bun.In(filter.Level))
	}

	if len(filter.Provider) > 0 {
		query = query.Where("provider in (?)", bun.In(filter.Provider))
	}

	if len(filter.Action) > 0 {
		query = query.Where("action in (?)", bun.In(filter.Action))
	}

	if filter.Path != "" {
		query = query.Where("path like LOWER(?)", "%"+filter.Path+"%")
	}

	if filter.Resource != "" {
		query = query.Where("resource like UPPER(?)", "%"+filter.Resource+"%")
	}

	if len(filter.TenantID) > 0 {
		allowedTenantsIds := filterAllowedTenants(userTenantsID, filter.TenantID)

		if len(allowedTenantsIds) == 0 {
			return model, 0, nil
		}

		conditions := buildQueryTenants(allowedTenantsIds, "OR")
		query = query.Where(conditions)
	}

	if len(filter.UserID) > 0 {
		query = query.Where("user_id in (?)", bun.In(filter.UserID))
	}

	if !filter.StartAt.IsZero() && !filter.EndAt.IsZero() {
		query = query.Where("created_at::TIMESTAMP BETWEEN TIMESTAMP ? AND TIMESTAMP ?", filter.StartAt, filter.EndAt)
	}

	conditions := buildQueryTenants(userTenantsID, "OR")
	query = query.Where(conditions)

	if err := query.Scan(ctx); err != nil {
		return nil, 0, err
	}

	count, _ := query.Count(ctx)

	return model, count, nil
}

func buildQueryTenants(tenants []int, operator string) string {

	if operator == "" {
		operator = "OR"
	}

	conditions := "("
	for i, id := range tenants {
		if i > 0 {
			conditions += fmt.Sprintf(" %s ", operator)
		}
		conditions += fmt.Sprintf("tenant_id @> '[%d]'", id)
	}
	conditions += ")"

	return conditions
}

func filterAllowedTenants(allowedTenants []int, requestedTenants []int) []int {

	allowedMap := make(map[int]bool)
	for _, id := range allowedTenants {
		allowedMap[id] = true
	}

	var res []int
	for _, id := range requestedTenants {
		if _, check := allowedMap[id]; check {
			res = append(res, id)
		}
	}

	return res
}
