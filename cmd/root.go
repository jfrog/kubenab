package cmd

import (
	"fmt"
	"os"

	"github.com/jfrog/kubenab/internal"
	"github.com/jfrog/kubenab/pkg/log"
	_log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubenab",
	Short: "K8s Image policy enforcer",
	// TODO: This 'long' description needs to be enhanced!
	Long: `Kubernetes Admission Webhook to enforce pulling of Docker images from the private registry.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		/// check if there are any config incompatibilities

		// print error if user set 'debug' flag but application was
		// compiled without debug abilities
		if internal.Debug && !internal.DebugAvail {
			log.Errorln("Debug output enabled but this application was not compiled with debug support.")
		} else if internal.Debug {
			// increase log level servity since this binary has
			// been compiled with the debug feature
			_log.SetLevel(_log.DebugLevel)
		}

		// catch if '--version' flag was set
		val, err := cmd.PersistentFlags().GetBool("version")
		if err != nil {
			return err
		}

		if val {
			log.Printf("Version......: %s\n", internal.Version)
			log.Printf("Build Date...: %s\n", internal.BuildDate)
			log.Printf("Commit.......: %s\n", internal.Commit)
			log.Printf("Debug capable: %t\n", internal.DebugAvail)

			os.Exit(0)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/kubenab/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&internal.Debug, "debug", false, "print Debug Messages (defaults to false)")
	rootCmd.PersistentFlags().Bool("version", false, "print version and build information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/kubenab")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
