package main

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"io/ioutil"
	"os"
	"time"
)

var config = Config{
	AnimationSpeed: 30,
	Fps:            60,
	Control: Control{
		Left:  pixelgl.KeyA,
		Right: pixelgl.KeyD,
		Up:    pixelgl.KeyW,
		Down:  pixelgl.KeyS,
	},
}

type Config struct {
	AnimationSpeed time.Duration
	Fps            time.Duration
	Control
}

func LoadConfiguration(file string, config *Config) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}

func SaveConfig() {
	js, _ := json.Marshal(config)
	ioutil.WriteFile("settings.json", js, 0644)
}
