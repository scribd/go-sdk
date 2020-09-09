package database

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
)

type ctxDatabaseMarker struct{}

type ctxDatabase struct {
	database *gorm.DB
}

var (
	ctxDatabaseKey = &ctxDatabaseMarker{}
)

// Extract takes the gorm.DB database from the context.
// If the ctxDatabase wasn't used, an error is returned.
func Extract(ctx context.Context) (*gorm.DB, error) {
	d, ok := ctx.Value(ctxDatabaseKey).(*ctxDatabase)
	if !ok || d == nil {
		return nil, fmt.Errorf("Unable to get the database")
	}

	return d.database, nil
}

// ToContext adds the gorm.DB database to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, db *gorm.DB) context.Context {
	d := &ctxDatabase{
		database: db,
	}
	return context.WithValue(ctx, ctxDatabaseKey, d)
}
