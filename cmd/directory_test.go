package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveDirSource(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "dir with spaces")
	require.NoError(t, os.Mkdir(dir, 0o700))
	file := filepath.Join(dir, "secret.txt")
	require.NoError(t, os.WriteFile(file, []byte("synthetic"), 0o600))
	unicodeDir := filepath.Join(root, "unicodé")
	require.NoError(t, os.Mkdir(unicodeDir, 0o700))

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "no arguments uses current directory", want: "."},
		{name: "empty argument uses current directory", args: []string{""}, want: "."},
		{name: "one directory", args: []string{dir}, want: dir},
		{name: "one file", args: []string{file}, want: file},
		{name: "unicode path", args: []string{unicodeDir}, want: unicodeDir},
		{name: "multiple arguments rejected", args: []string{file, file}, wantErr: true},
		{name: "nonexistent target rejected", args: []string{filepath.Join(root, "missing")}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolveDirSource(test.args)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestDirectoryCommandArgs(t *testing.T) {
	require.NoError(t, directoryCmd.Args(directoryCmd, nil))
	require.NoError(t, directoryCmd.Args(directoryCmd, []string{"target"}))
	require.Error(t, directoryCmd.Args(directoryCmd, []string{"first", "second"}))
}
