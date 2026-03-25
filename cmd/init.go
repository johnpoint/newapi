package cmd

import (
	"github.com/spf13/cobra"

	"newapi/internal/generator"
)

var initCmd = &cobra.Command{
	Use:   "init [module-name]",
	Short: "Initialize a new Go API project",
	Long: `Initialize a new Go API project in the current directory.

Creates the standard directory layout:
  cmd/server/main.go        - server entry point
  internal/endpoints/       - generated endpoint files
  internal/services/        - generated service files
  internal/schemas/         - generated schema files

If a module name is provided, go mod init is run automatically.
If go.mod already exists, the module name is read from it.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := ""
		if len(args) == 1 {
			moduleName = args[0]
		}
		if err := generator.InitProject(moduleName); err != nil {
			fatal("%v", err)
		}
	},
}
