package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"rashomon/cmd"
	"rashomon/conf"
	"rashomon/pkg/logger"
	"rashomon/pkg/mysql"
	"rashomon/pkg/redis"
)

func main() {
	// 应用的主入口，默认调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   "rashomon",
		Short: `use "-h" flag to see all subcommands`,
		Long:  `use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			conf.Init()
			logger.Init()
			redis.Init()
			mysql.Init()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			logger.Quit()
		},
	}
	// 注册子命令,用于有一些命令需要编写时候
	rootCmd.AddCommand(
		cmd.WebServer,
	)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed to run app with %v: %s\n", os.Args, err.Error())
	}
}
