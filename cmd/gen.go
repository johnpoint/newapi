package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"newapi/internal/generator"
)

var configFile string

var genApiCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成 api",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, err := generator.ReadModuleName()
		if err != nil {
			fatal("%v", err)
		}

		config, err := os.ReadFile(configFile)
		if err != nil {
			fatal("读取配置文件失败: %v", err)
		}

		apiDesc, err := generator.ParseConfig(config)
		if err != nil {
			fatal("解析配置文件失败: %v", err)
		}

		var total int64
		for _, v := range apiDesc.Group {
			for _, api := range v.Apis {
				a := generator.ApiTemp{
					ProjectName: projectName,
					ApiName:     api.ApiName,
					Method:      strings.ToUpper(api.Method),
					RouterPath:  strings.ReplaceAll(strings.ReplaceAll(api.Path, "{", ":"), "}", ""),
					Path:        api.Path,
					GenCommand:  strings.Join(os.Args, " "),
					Desc:        api.Desc,
					Tag:         v.Name,
					PathVar:     generator.GetPathVar(api.Path),
				}
				fmt.Printf("正在生成: %s\t%s\t%s\t%s\n", a.ApiName, a.Method, a.Path, a.Desc)
				if err := generator.GenEndpoint(a); err != nil {
					fatal("%v", err)
				}
				if err := generator.GenService(a); err != nil {
					fatal("%v", err)
				}
				if err := generator.GenSchema(a); err != nil {
					fatal("%v", err)
				}
				total++
			}
		}
		fmt.Printf("生成完成，共 %d 个接口\n", total)
	},
}
