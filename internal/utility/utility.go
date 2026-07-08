package utility

import (
	"os"
	"path/filepath"
)

func GetSocketPath(socketAddr string) string {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		dir = "/tmp" // Fallback
	}
	return filepath.Join(dir, socketAddr)
}
