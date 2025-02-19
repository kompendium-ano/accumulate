package node_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/AccumulateNetwork/accumulate/config"
	"github.com/AccumulateNetwork/accumulate/internal/logging"
	"github.com/AccumulateNetwork/accumulate/internal/node"
	acctesting "github.com/AccumulateNetwork/accumulate/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeLifecycle(t *testing.T) {
	acctesting.SkipPlatformCI(t, "darwin", "requires setting up localhost aliases")

	// Configure
	opts := acctesting.NodeInitOptsForLocalNetwork(t.Name(), acctesting.GetIP())
	opts.WorkDir = t.TempDir()
	opts.Logger = logging.NewTestLogger(t, "plain", config.DefaultLogLevels, false)
	require.NoError(t, node.Init(opts))

	// Start
	nodeDir := filepath.Join(opts.WorkDir, "Node0")
	daemon, err := acctesting.RunDaemon(acctesting.DaemonOptions{
		Dir:       nodeDir,
		LogWriter: logging.TestLogWriter(t),
	}, t.Cleanup)
	require.NoError(t, err)
	require.NoError(t, daemon.Stop())

	// procfs is a linux thing
	if runtime.GOOS != "linux" {
		return
	}

	fds := filepath.Join("/proc", fmt.Sprint(os.Getpid()), "fd")
	entries, err := os.ReadDir(fds)
	require.NoError(t, err)
	for _, e := range entries {
		if e.Type()&os.ModeSymlink == 0 {
			continue
		}

		file, err := filepath.EvalSymlinks(filepath.Join(fds, e.Name()))
		if err != nil {
			continue
		}

		rel, err := filepath.Rel(nodeDir, file)
		require.NoError(t, err)

		if strings.HasPrefix(rel, "../") {
			continue
		}

		assert.Failf(t, "Files are still open after the node was shut down", "%q is open", rel)
	}
}
