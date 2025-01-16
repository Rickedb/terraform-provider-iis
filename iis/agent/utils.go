package agent

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func generateHash(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func toPascalCase(value bool) string {
	bVal := "False"
	if value {
		bVal = "True"
	}

	return bVal
}

func appendEnsurePhysicalPath(sb *strings.Builder, path string) {
	physicalPath := strings.ReplaceAll(path, "/", `\`)
	sb.WriteString(fmt.Sprintf(`
		$path='%v'
		if (!(Test-Path $path)){
            New-Item -ItemType Directory -Path $path;
        }
	`, physicalPath))
}
