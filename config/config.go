package config

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"io/ioutil"
	"os"
	"time"
	"github.com/rs/xid"
)

var ID = xid.New().String()

var Conf = Config{
	AnimationSpeed:    time.Second / 30,
	ServerHandleSpeed: time.Second / 200,
	Fps:               time.Second / 300,
	LeftKey:           pixelgl.KeyA,
	RightKey:          pixelgl.KeyD,
	UpKey:             pixelgl.KeyW,
	DownKey:           pixelgl.KeyS,
	ShootKey:          pixelgl.KeySpace,
	MeleeKey:          pixelgl.KeyLeftControl,
	KnifeKey:          pixelgl.Key1,
	HandgunKey:        pixelgl.Key2,
	RifleKey:          pixelgl.Key3,
	ShotgunKey:        pixelgl.Key4,
	ReloadKey:         pixelgl.KeyR,
	BarrelKey:         pixelgl.KeyB,
	Online:            false,
	LoungeUri:         "tcp://localhost:31415/lounge",
}

type Config struct {
	AnimationSpeed    time.Duration
	Fps               time.Duration
	ServerHandleSpeed time.Duration
	LeftKey           pixelgl.Button
	RightKey          pixelgl.Button
	UpKey             pixelgl.Button
	DownKey           pixelgl.Button
	ShootKey          pixelgl.Button
	MeleeKey          pixelgl.Button
	KnifeKey          pixelgl.Button
	HandgunKey        pixelgl.Button
	RifleKey          pixelgl.Button
	ShotgunKey        pixelgl.Button
	ReloadKey         pixelgl.Button
	BarrelKey         pixelgl.Button
	Online            bool
	LoungeUri         string
}

func LoadJson(file string, config interface{}) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}

func SaveConfig(file string) {
	return
	js, _ := json.Marshal(Conf)
	ioutil.WriteFile(file, js, 0644) // todo: find better way to save settings file
}
