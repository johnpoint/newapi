package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd    = &cobra.Command{}
	configPath string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func init() {
	newApiCommand.Flags().StringVarP(&a.ApiName, "name", "n", "TestApi", "api 名称，使用大驼峰")
	newApiCommand.Flags().StringVarP(&a.Desc, "desc", "d", "auto generate api", "api 描述")
	newApiCommand.Flags().StringVarP(&a.Tag, "tag", "t", "test", "api 标签")
	newApiCommand.Flags().StringVarP(&a.Method, "method", "m", "get", "请求方法")
	newApiCommand.Flags().StringVarP(&a.Path, "path", "p", "/ping", "请求路径")

	genApiCommand.Flags().StringVarP(&file, "config", "c", "", "声明文件")

	rootCmd.AddCommand(newApiCommand)
	rootCmd.AddCommand(genApiCommand)
}
