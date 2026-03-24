package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed assets/endpoint.tmpl
var endpoint string

//go:embed assets/service.tmpl
var service string

//go:embed assets/schema.tmpl
var schema string

type ApiTemp struct {
	ProjectName string
	ApiName     string
	Path        string
	RouterPath  string
	Desc        string
	Tag         string
	Method      string
	GenCommand  string
	PathVar     []string
}

var newApiCommand = &cobra.Command{
	Use:   "new",
	Short: "新建 api",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		desc, _ := cmd.Flags().GetString("desc")
		tag, _ := cmd.Flags().GetString("tag")
		method, _ := cmd.Flags().GetString("method")
		path, _ := cmd.Flags().GetString("path")

		projectName, err := readModuleName()
		if err != nil {
			fatal("%v", err)
		}

		a := ApiTemp{
			ProjectName: projectName,
			ApiName:     name,
			Desc:        desc,
			Tag:         tag,
			Method:      strings.ToUpper(method),
			Path:        path,
			RouterPath:  strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", ""),
			GenCommand:  strings.Join(os.Args, " "),
			PathVar:     GetPathVar(path),
		}

		if err := genEndpoint(a); err != nil {
			fatal("%v", err)
		}
		if err := genService(a); err != nil {
			fatal("%v", err)
		}
		if err := genSchema(a); err != nil {
			fatal("%v", err)
		}
	},
}

func genEndpoint(a ApiTemp) error {
	path := fmt.Sprintf("endpoints/%s_ep.go", CamelCaseToUnderscore(a.ApiName))
	return generateFile(path, endpoint, a, false)
}

func genService(a ApiTemp) error {
	path := fmt.Sprintf("services/%s_srv.go", CamelCaseToUnderscore(a.ApiName))
	return generateFile(path, service, a, false)
}

func genSchema(a ApiTemp) error {
	path := fmt.Sprintf("schemas/%s_schema.go", CamelCaseToUnderscore(a.ApiName))
	return generateFile(path, schema, a, true)
}
