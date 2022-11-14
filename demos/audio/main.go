package main

import (
	"sckorok"
	"sckorok/asset"
	"sckorok/audio"
	"sckorok/game"
	"sckorok/gui"

	"fmt"
	"log"
	"sckorok/gfx/font"
)

func main() {
	fmt.Println("Hello Audio!!")
	options := sckorok.Options{
		Width:  320,
		Height: 480,
		Title:  "Audio Test",
	}
	sckorok.Run(&options, &MainScene{})
}

type MainScene struct {
	wav uint16
	ogg uint16
}

func (*MainScene) Load() {
	asset.Font.LoadTrueType("ttf", "OCRAEXT.TTF", font.ASCII(24))

	asset.Audio.Load("birds.wav", true)
	asset.Audio.Load("ambient.ogg", true)
}

func (m *MainScene) OnEnter(g *game.Game) {
	font, _ := asset.Font.Get("ttf")
	gui.SetFont(font)
	gui.SetVirtualResolution(320, 0)

	m.wav, _ = asset.Audio.Get("birds.wav")
	m.ogg, _ = asset.Audio.Get("ambient.ogg")
}

func (m *MainScene) Update(dt float32) {
	if gui.Button(1, gui.Rect{X: 100, Y: 100, W: 0, H: 0}, "Play", nil).JustPressed() {
		audio.PlayMusic(m.ogg)
		log.Println("play audio")
	}

	if gui.Button(2, gui.Rect{X: 100, Y: 140, W: 0, H: 0}, "Stop", nil).JustPressed() {
		// stop audio
		audio.StopMusic()
		log.Println("stop audio")
	}

	if gui.Button(3, gui.Rect{X: 180, Y: 100, W: 0, H: 0}, "Pause", nil).JustPressed() {
		audio.PauseMusic()
		log.Println("pause audio")
	}

	if gui.Button(4, gui.Rect{X: 180, Y: 140, W: 0, H: 0}, "Resume", nil).JustPressed() {
		audio.ResumeMusic()
		log.Println("resume audio")
	}
}

func (m *MainScene) OnExit() {
	audio.Destroy()
}
