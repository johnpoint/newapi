package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var configFile string

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
		projectName, err := readModuleName()
		if err != nil {
			fatal("%v", err)
		}

		config, err := os.ReadFile(configFile)
		if err != nil {
			fatal("读取配置文件失败: %v", err)
		}

		var apiDesc ApiDesc
		if err := yaml.Unmarshal(config, &apiDesc); err != nil {
			fatal("解析配置文件失败: %v", err)
		}

		var total int64
		for _, v := range apiDesc.Group {
			for _, api := range v.Apis {
				apiTemp := ApiTemp{
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
				if err := genEndpoint(apiTemp); err != nil {
					fatal("%v", err)
				}
				if err := genService(apiTemp); err != nil {
					fatal("%v", err)
				}
				if err := genSchema(apiTemp); err != nil {
					fatal("%v", err)
				}
				total++
			}
		}
		fmt.Printf("生成完成，共 %d 个接口\n", total)
	},
}
