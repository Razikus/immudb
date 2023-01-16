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

package sql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConditionalRowReader(t *testing.T) {
	dummyr := &dummyRowReader{failReturningColumns: true}

	rowReader := newConditionalRowReader(dummyr, &Bool{val: true})

	_, err := rowReader.Columns(context.Background())
	require.Equal(t, errDummy, err)

	err = rowReader.InferParameters(context.Background(), nil)
	require.Equal(t, errDummy, err)

	dummyr.failInferringParams = true

	err = rowReader.InferParameters(context.Background(), nil)
	require.Equal(t, errDummy, err)
}
