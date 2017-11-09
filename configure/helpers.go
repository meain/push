package configure

import (
	"os/user"
	"strings"
)

func expandFilePath(filePath string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return strings.Replace(filePath, "~", homeDir, 1)
}
