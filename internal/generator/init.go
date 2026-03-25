package generator

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed templates/main.tmpl
var mainTmpl string

// ProjectConfig holds the configuration used to initialize a new project.
type ProjectConfig struct {
	ModuleName string
}

// InitProject scaffolds a new Go API project in the current directory.
//
// It creates the standard directory layout, generates cmd/server/main.go,
// and optionally runs `go mod init` if no go.mod is present.
//
// If moduleName is empty, the module name is read from an existing go.mod.
func InitProject(moduleName string) error {
	// Resolve module name.
	if moduleName == "" {
		name, err := ReadModuleName()
		if err != nil {
			return fmt.Errorf("无法确定模块名：当前目录不含 go.mod，请提供模块名作为参数（例如：newapi init myapp）")
		}
		moduleName = name
	}

	// Create directory layout.
	dirs := []string{
		"cmd/server",
		"internal/endpoints",
		"internal/services",
		"internal/schemas",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %w", d, err)
		}
		fmt.Printf("  创建目录  %s/\n", d)
	}

	// Initialize go.mod if not present.
	if exists, _ := pathExists("go.mod"); !exists {
		fmt.Printf("  运行      go mod init %s\n", moduleName)
		cmd := exec.Command("go", "mod", "init", moduleName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("go mod init 失败: %w", err)
		}
	}

	// Generate cmd/server/main.go.
	mainPath := "cmd/server/main.go"
	if exists, _ := pathExists(mainPath); exists {
		fmt.Printf("  跳过      %s (已存在)\n", mainPath)
	} else {
		if err := renderToFile(mainPath, mainTmpl, ProjectConfig{ModuleName: moduleName}); err != nil {
			return err
		}
		fmt.Printf("  生成      %s\n", mainPath)
	}

	fmt.Printf("\n项目初始化完成！模块名: %s\n", moduleName)
	fmt.Println("\n下一步：")
	fmt.Printf("  newapi new -n Ping -m get -p /ping\n")
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run ./cmd/server\n")
	return nil
}

// renderToFile renders a text/template string to a new file.
func renderToFile(path, tmplStr string, data any) error {
	if err := os.MkdirAll(dirOf(path), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("创建文件失败 %s: %w", path, err)
	}
	defer f.Close()

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染模板失败: %w", err)
	}
	_, err = f.WriteString(buf.String())
	return err
}
