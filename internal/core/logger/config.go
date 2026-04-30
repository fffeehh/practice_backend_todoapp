package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level string `envcongig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER" requires:"true"`
}

func NewConfig() (Config, error) {
	var config Config
	
	// первый аргумент - это преффикс. Он автоматически подставляется в начало названия поля в envcongig.
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil 
}

// конструктор, который выбрасывает панику, если переменные не удалось проинициализировать. Нужен тогда, когда нам нет смысла как-то обрабатывать ошибки, а нужно сразу положить программу
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}
	return config
}
