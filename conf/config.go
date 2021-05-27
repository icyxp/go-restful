package conf

import (
	"fmt"
	"strings"

	"go-restful/lib/log"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// InitDir init config
func Init(confPath string) error {
	err := initConfig(confPath)
	if err != nil {
		return err
	}
	return nil
}

func initConfig(confPath string) error {
	if confPath != "" {
		viper.SetConfigFile(confPath)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("xp")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return errors.WithStack(err)
	}

	watchConfig()

	return nil
}

// Monitor configuration file changes and hot load programs
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

// init log
func InitLog() {
	config := log.Config{
		Writers:          viper.GetString("logout.writers"),
		LoggerLevel:      viper.GetString("logout.logger_level"),
		LoggerFile:       viper.GetString("logout.logger_file"),
		LoggerWarnFile:   viper.GetString("logout.logger_warn_file"),
		LoggerErrorFile:  viper.GetString("logout.logger_error_file"),
		LogFormatText:    viper.GetBool("logout.log_format_text"),
		LogRollingPolicy: viper.GetString("logout.log_rolling_policy"),
		LogRotateDate:    viper.GetInt("logout.log_rotate_date"),
		LogRotateSize:    viper.GetInt("logout.log_rotate_size"),
		LogBackupCount:   viper.GetInt("logout.log_backup_count"),
	}

	err := log.NewLogger(&config, log.InstanceZapLogger)
	if err != nil {
		fmt.Printf("InitWithConfig err: %v", err)
	}
}
