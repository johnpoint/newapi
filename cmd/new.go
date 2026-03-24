package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"newapi/internal/generator"
)

var newApiCommand = &cobra.Command{
	Use:   "new",
	Short: "新建 api",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("desc")
		tag, _ := cmd.Flags().GetString("tag")
		method, _ := cmd.Flags().GetString("method")
		path, _ := cmd.Flags().GetString("path")

		projectName, err := generator.ReadModuleName()
		if err != nil {
			fatal("%v", err)
		}

		a := generator.ApiTemp{
			ProjectName: projectName,
			ApiName:     name,
			Desc:        desc,
			Tag:         tag,
			Method:      strings.ToUpper(method),
			Path:        path,
			RouterPath:  strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", ""),
			GenCommand:  strings.Join(os.Args, " "),
			PathVar:     generator.GetPathVar(path),
		}

		if err := generator.GenEndpoint(a); err != nil {
			fatal("%v", err)
		}
		if err := generator.GenService(a); err != nil {
			fatal("%v", err)
		}
		if err := generator.GenSchema(a); err != nil {
			fatal("%v", err)
		}
		fmt.Printf("已生成: %s %s %s\n", a.Method, a.Path, a.ApiName)
	},
}
