/*
Copyright 2022 Codenotary Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/codenotary/immudb/pkg/client"
	"github.com/codenotary/immudb/pkg/client/errors"
	"github.com/stretchr/testify/require"
)

func TestGRPCError(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_LEVEL")

	bs, cli, _ := setupTestServerAndClientWithToken(t)

	t.Run("errors with token-based auth", func(t *testing.T) {

		_, err := cli.Login(context.Background(), []byte(`immudb`), []byte(`wrong`))

		require.Equal(t, err.(errors.ImmuError).Error(), "invalid user name or password")
		require.Equal(t, err.(errors.ImmuError).Cause(), "crypto/bcrypt: hashedPassword is not the hash of the given password")
		require.Equal(t, err.(errors.ImmuError).Code(), errors.CodSqlserverRejectedEstablishmentOfSqlconnection)
		require.Equal(t, int32(0), err.(errors.ImmuError).RetryDelay())
		require.NotNil(t, err.(errors.ImmuError).Stack())

	})

	t.Run("errors with session-based auth", func(t *testing.T) {
		cli := bs.NewClient(client.DefaultOptions())

		err := cli.OpenSession(context.Background(), []byte(`immudb`), []byte(`wrong`), "defaultdb")

		require.Equal(t, err.(errors.ImmuError).Error(), "invalid user name or password")
		require.Equal(t, err.(errors.ImmuError).Cause(), "crypto/bcrypt: hashedPassword is not the hash of the given password")
		require.Equal(t, err.(errors.ImmuError).Code(), errors.CodSqlserverRejectedEstablishmentOfSqlconnection)
		require.Equal(t, int32(0), err.(errors.ImmuError).RetryDelay())
		require.NotNil(t, err.(errors.ImmuError).Stack())
	})
}
