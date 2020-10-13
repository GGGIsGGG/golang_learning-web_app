package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Print("viper.ReadInConfig() failed, err: %V\n", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件修改了！")
	})
	return
}
