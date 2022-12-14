package diskutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spiffe/spire/test/spiretest"
	"github.com/stretchr/testify/require"
)

func TestAtomicWriteFile(t *testing.T) {
	dir := spiretest.TempDir(t)

	tests := []struct {
		name              string
		data              []byte
		mode              os.FileMode
		expectPosixMode   os.FileMode
		expectWindowsMode os.FileMode
	}{
		{
			name:              "basic test",
			data:              []byte("Hello, World"),
			mode:              0600,
			expectPosixMode:   0600,
			expectWindowsMode: 0666,
		},
		{
			name:              "empty",
			data:              []byte{},
			mode:              0440,
			expectPosixMode:   0440,
			expectWindowsMode: 0444,
		},
		{
			name:              "binary",
			data:              []byte{0xFF, 0, 0xFF, 0x3D, 0xD8, 0xA9, 0xDC, 0xF0, 0x9F, 0x92, 0xA9},
			mode:              0644,
			expectPosixMode:   0644,
			expectWindowsMode: 0666,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(dir, "file")
			err := AtomicWriteFile(file, tt.data, tt.mode)
			require.NoError(t, err)

			info, err := os.Stat(file)
			require.NoError(t, err)
			switch runtime.GOOS {
			case "windows":
				require.EqualValues(t, tt.expectWindowsMode, info.Mode())
			default:
				require.EqualValues(t, tt.expectPosixMode, info.Mode())
			}

			content, err := os.ReadFile(file)
			require.NoError(t, err)
			require.Equal(t, tt.data, content)

			require.NoError(t, os.Remove(file))
		})
	}
}
