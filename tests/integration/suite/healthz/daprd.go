/*
Copyright 2023 The Dapr Authors
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

package healthz

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapr/dapr/tests/integration/framework"
	procdaprd "github.com/dapr/dapr/tests/integration/framework/process/daprd"
	"github.com/dapr/dapr/tests/integration/suite"
)

func init() {
	suite.Register(new(daprd))
}

// daprd tests that Dapr responds to healthz requests.
type daprd struct {
	proc *procdaprd.Daprd
}

func (d *daprd) Setup(t *testing.T) []framework.Option {
	d.proc = procdaprd.New(t)
	return []framework.Option{
		framework.WithProcesses(d.proc),
	}
}

func (d *daprd) Run(t *testing.T, ctx context.Context) {
	d.proc.WaitUntilRunning(t, ctx)

	reqURL := fmt.Sprintf("http://localhost:%d/v1.0/healthz", d.proc.PublicPort())

	assert.Eventually(t, func() bool {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())
		return resp.StatusCode == http.StatusNoContent
	}, time.Second*10, 100*time.Millisecond)
}
