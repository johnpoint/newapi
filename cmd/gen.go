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

		var total int64
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
					PathVar:     GetPathVar(api.Path),
				}
				fmt.Printf("正在生成: %s\t%s\t%s\t%s\n", apiTemp.ApiName, apiTemp.Method, apiTemp.Path, apiTemp.Desc)
				genEndpoint(apiTemp)
				genService(apiTemp)
				genSchema(apiTemp)
				total++
			}
		}
		fmt.Printf("生成完成，共 %d 个接口\n", total)
	},
}

func GetPathVar(path string) []string {
	var paths []string
	var tmpPathVar string
	var catch bool
	for _, v := range path {
		if catch && v == '}' {
			catch = false
			paths = append(paths, tmpPathVar)
			tmpPathVar = ""
			continue
		}
		if v == '{' {
			catch = true
			continue
		}
		if catch {
			tmpPathVar += string(v)
		}
	}
	return paths
}
