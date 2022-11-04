package main

import (
	korok "sckorok"
	"sckorok/asset"
	"sckorok/game"
	"sckorok/gfx/font"
	"sckorok/math/f32"
)

type MainScene struct {
}

func (*MainScene) Load() {
	asset.Font.LoadBitmap("font1", "font.png", "font.json")
	asset.Font.LoadTrueType("font2", "OCRAEXT.TTF", font.ASCII(24))
}

func (*MainScene) OnEnter(g *game.Game) {
	font, _ := asset.Font.Get("font1")

	// show "Hello world"
	entity := korok.Entity.New()

	text := korok.Text.NewComp(entity)
	text.SetFont(font)
	text.SetText("Hello Korok!")

	xf := korok.Transform.NewComp(entity)
	xf.SetPosition(f32.Vec2{240, 160})
}

func (*MainScene) Update(dt float32) {

}

func (*MainScene) OnExit() {
}

func main() {
	options := korok.Options{
		Title:  "Text Rendering",
		Width:  480,
		Height: 320,
	}
	korok.Run(&options, &MainScene{})
}
