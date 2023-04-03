package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var file string

type ApiDesc struct {
	Group []ApiGroup `yaml:"group"`
}

type ApiGroup struct {
	Name string `yaml:"name"`
	Apis []Api  `yaml:"apis"`
}

type Api struct {
	ApiName string `yaml:"name"`
	Path    string `yaml:"path"`
	Desc    string `yaml:"desc"`
	Method  string `yaml:"method"`
}

var genApiCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成 api",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open("go.mod")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		all, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}

		var projectName string

		s := strings.Split(string(all), "\n")
		for _, v := range s {
			vv := strings.Split(v, " ")
			if len(vv) == 2 && vv[0] == "module" {
				projectName = vv[1]
				break
			}
		}

		cf, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer cf.Close()
		config, err := io.ReadAll(cf)
		if err != nil {
			panic(err)
		}
		var apiDesc ApiDesc
		err = yaml.Unmarshal(config, &apiDesc)
		if err != nil {
			panic(err)
		}

		for _, v := range apiDesc.Group {
			for _, api := range v.Apis {
				var apiTemp = ApiTemp{
					ProjectName: projectName,
					ApiName:     api.ApiName,
					Method:      strings.ToUpper(api.Method),
					RouterPath:  strings.ReplaceAll(strings.ReplaceAll(api.Path, "{", ":"), "}", ""),
					Path:        api.Path,
					GenCommand:  strings.Join(os.Args, " "),
					Desc:        api.Desc,
					Tag:         v.Name,
				}
				fmt.Printf("gen: %+v\n", apiTemp)
				genEndpoint(apiTemp)
				genService(apiTemp)
				genSchema(apiTemp)
			}
		}
	},
}
