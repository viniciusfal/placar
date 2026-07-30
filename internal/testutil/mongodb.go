// internal/testutil/mongodb.go
package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	tcmongodb "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/viniciusfal/placar/internal/platform/mongodb"
)

func SetupTestDB(t *testing.T) *mongo.Database {
	t.Helper()

	ctx := context.Background()
	container, err := tcmongodb.Run(ctx, "mongo:7")
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("falha ao encerrar container: %v", err)
		}
	})

	uri, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongodb.Connect(ctx, uri)
	require.NoError(t, err)

	return client.Database("test")
}
