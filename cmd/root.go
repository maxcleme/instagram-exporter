/*
Copyright © 2020 maxcleme <maximeclement93+git@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/maxcleme/instagram-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "instagram-exporter",
	Short: "Prometheus exporter for Instagram account metrics",
	Long: `Prometheus exporter for Instagram account metrics

It exports the following metrics :
- instagram_media_like_total : Total likes count by [username, media_id]
- instagram_media_comment_total : Total comments count by [username, media_id]
- instagram_media_total : Total media count by [username]
- instagram_follower_total : Total followers count by [username]
- instagram_following_total : Total following count by [username]
`,
	Run: func(cmd *cobra.Command, args []string) {
		level, _ := cmd.Flags().GetString("log_level")
		lv, err := logrus.ParseLevel(level)
		if err != nil {
			logrus.WithError(err).Fatal("cannot parse log level")
		}
		logrus.SetLevel(lv)

		usernames, _ := cmd.Flags().GetStringSlice("ig_target_users")
		login, _ := cmd.Flags().GetString("ig_login")
		password, _ := cmd.Flags().GetString("ig_password")
		tokenPath, _ := cmd.Flags().GetString("ig_token_path")
		cacheDuration, _ := cmd.Flags().GetDuration("cache_duration")

		c, err := collector.Instagram(
			collector.WithTokenPath(tokenPath),
			collector.WithLogin(login, password),
			collector.WithTargets(usernames...),
			collector.WithCacheDuration(cacheDuration),
		)
		if err != nil {
			logrus.WithError(err).Fatal("cannot create instagram collector")
		}

		if err := prometheus.Register(c); err != nil {
			logrus.WithError(err).Fatal("cannot register prom exporter")
		}

		path, _ := cmd.Flags().GetString("http_path")
		port, _ := cmd.Flags().GetInt("http_port")

		http.Handle(path, promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			logrus.WithError(err).Fatal("exporter stop abnormally")
		}

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.instagram-exporter.yaml)")

	// misc flags
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().String("log_level", "info", "Log level")
	rootCmd.Flags().Int("http_port", 2112, "Port used by Prometheus to expose metrics")
	rootCmd.Flags().String("http_path", "/metrics", "Path used by Prometheus to expose metrics")
	rootCmd.Flags().Duration("cache_duration", 5*time.Minute, "Cache duration")

	// exporter flags
	rootCmd.Flags().String("ig_login", "", "Instagram login")
	rootCmd.Flags().String("ig_password", "", "Instagram password")

	rootCmd.Flags().StringSlice("ig_target_users", nil, "Instagram usernames to fetch metrics")
	rootCmd.MarkFlagRequired("ig_target_users")

	rootCmd.Flags().String("ig_token_path", filepath.Join(os.TempDir(), "instagram-exporter", ".ig_token"), "Instagram token file location")
	rootCmd.MarkFlagRequired("ig_token_path")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".instagram-exporter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".instagram-exporter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// workaround to set cobra required flags with viper
	postInitCommands([]*cobra.Command{rootCmd})
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
