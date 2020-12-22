package client

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/server"
	"github.com/codenotary/immudb/pkg/server/servertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestDatabasesSwitching(t *testing.T) {
	options := server.DefaultOptions().WithAuth(true)
	bs := servertest.NewBufconnServer(options)

	go func() {
		bs.Start()
	}()

	defer os.RemoveAll(options.Dir)

	ctx := context.Background()

	err := client.CreateDatabase(ctx, &schema.Database{
		Databasename: "db1",
	})
	require.Nil(t, err)
	resp, err := client.UseDatabase(ctx, &schema.Database{
		Databasename: "db1",
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Token)
	_, err = client.VerifiedSet(ctx, []byte(`db1-my`), []byte(`item`))
	assert.Nil(t, err)

	err = client.CreateDatabase(ctx, &schema.Database{
		Databasename: "db2",
	})
	assert.Nil(t, err)
	resp2, err := client.UseDatabase(ctx, &schema.Database{
		Databasename: "db2",
	})

	md := metadata.Pairs("authorization", resp2.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Token)
	_, err = client.VerifiedSet(ctx, []byte(`db2-my`), []byte(`item`))
	assert.Nil(t, err)

	vi, err := client.VerifiedGet(ctx, []byte(`db1-my`))
	assert.Error(t, err)
	assert.Nil(t, vi)
}

type PasswordReader struct {
	Pass       []string
	callNumber int
}

func (pr *PasswordReader) Read(msg string) ([]byte, error) {
	if len(pr.Pass) <= pr.callNumber {
		log.Fatal("Application requested the password more times than number of passwords supplied")
	}
	pass := []byte(pr.Pass[pr.callNumber])
	pr.callNumber++
	return pass, nil
}
