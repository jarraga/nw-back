package seed

import (
	"context"

	"nw-back/internal/postgres"
	"nw-back/internal/resetdb"
)

func resetDatabase(ctx context.Context) error {
	return resetdb.Reset(ctx, postgres.DB)
}
