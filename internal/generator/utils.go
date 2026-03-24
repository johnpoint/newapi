package generator

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// ReadModuleName reads the Go module name from go.mod in the current directory.
func ReadModuleName() (string, error) {
	f, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("打开 go.mod 失败: %w", err)
	}
	defer f.Close()

	all, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("读取 go.mod 失败: %w", err)
	}

	for _, line := range strings.Split(string(all), "\n") {
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) == 2 && parts[0] == "module" {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf("go.mod 中未找到 module 声明")
}

// CamelCaseToUnderscore converts PascalCase/camelCase to snake_case.
// e.g. GetUser → get_user, CreateOrderItem → create_order_item
func CamelCaseToUnderscore(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i == 0 {
			b.WriteRune(unicode.ToLower(r))
		} else {
			if unicode.IsUpper(r) {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

// GetPathVar extracts path parameter names from a path like /user/{id}/posts/{postId}.
func GetPathVar(path string) []string {
	var vars []string
	var cur strings.Builder
	inParam := false
	for _, r := range path {
		switch {
		case r == '{':
			inParam = true
		case r == '}' && inParam:
			inParam = false
			vars = append(vars, cur.String())
			cur.Reset()
		case inParam:
			cur.WriteRune(r)
		}
	}
	return vars
}

// pathExists returns (true, nil) if path exists, (false, nil) if it does not,
// or (false, err) for any other OS error (e.g. permission denied).
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("访问路径失败 %s: %w", path, err)
}

// dirOf returns the directory component of a slash- or backslash-separated path.
func dirOf(path string) string {
	idx := strings.LastIndexAny(path, "/\\")
	if idx < 0 {
		return "."
	}
	return path[:idx]
}
