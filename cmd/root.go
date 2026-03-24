package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "newapi",
	Short: "Go API 样板代码生成工具",
	Long: `newapi 是一个用于快速生成 Go API 样板代码的 CLI 工具。
基于 gin-gonic 和 go-bootstrap 框架，自动生成 endpoint、service、schema 三层文件。`,
}

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
	newApiCommand.Flags().StringP("name", "n", "TestApi", "api 名称，使用大驼峰")
	newApiCommand.Flags().StringP("desc", "d", "auto generate api", "api 描述")
	newApiCommand.Flags().StringP("tag", "t", "test", "api 标签")
	newApiCommand.Flags().StringP("method", "m", "get", "请求方法")
	newApiCommand.Flags().StringP("path", "p", "/ping", "请求路径")

	genApiCommand.Flags().StringVarP(&configFile, "config", "c", "", "声明文件")
	_ = genApiCommand.MarkFlagRequired("config")

	rootCmd.AddCommand(newApiCommand)
	rootCmd.AddCommand(genApiCommand)
}
