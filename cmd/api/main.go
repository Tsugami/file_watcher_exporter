package main

import (
	"file_watcher_exporter/internal/server"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("file_watcher_exporter")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.file_watcher_exporter")
	viper.AddConfigPath("/etc/file_watcher_exporter")
	viper.SetConfigType("yaml")

	viper.SetDefault("port", 5428)
	viper.SetDefault("host", "localhost")
	viper.SetDefault("dir", "/tmp")
	viper.SetDefault("exclude", []string{".*"})
	viper.SetDefault("include", []string{"*.*"})

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	config := &server.Config{
		Port:    viper.GetInt("port"),
		Host:    viper.GetString("host"),
		Dir:     viper.GetString("dir"),
		Exclude: viper.GetStringSlice("exclude"),
		Include: viper.GetStringSlice("include"),
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		config.SetDir(viper.GetString("dir"))
		config.SetExclude(viper.GetStringSlice("exclude"))
		config.SetInclude(viper.GetStringSlice("include"))
	})

	viper.WatchConfig()

	fmt.Printf("config loaded: %+v\n", config)
	server := server.NewServer(config)

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

	fmt.Printf("server started on port %d\n", config.Port)
}
