package main

import (
	"sckorok"
	"sckorok/asset"
	"sckorok/game"
	"sckorok/gfx"
	"sckorok/math/f32"
)

var s_Vertices = []gfx.PosTexColorVertex{
	{X: 100.0, Y: 100.0, U: 1, V: 1, RGBA: 0xffffffff},
	{X: -100.0, Y: -100.0, U: 0, V: 0, RGBA: 0xffffffff},
	{X: 100.0, Y: -100.0, U: 1, V: 0, RGBA: 0xffffffff},
	{X: -100.0, Y: 100.0, U: 0, V: 1, RGBA: 0xffffffff},
}

var s_Index = []uint16{
	3, 1, 2,
	3, 2, 0,
}

type MainScene struct {
}

func (*MainScene) Load() {
	asset.Texture.Load("face.png")
}

func (*MainScene) OnEnter(g *game.Game) {
	tex2d := asset.Texture.Get("face.png")
	// show mesh comp
	entity := sckorok.Entity.New()

	comp := sckorok.Mesh.NewComp(entity)
	mesh := &comp.Mesh

	mesh.SetIndex(s_Index)
	mesh.SetVertex(s_Vertices)
	mesh.Setup()
	mesh.SetTexture(tex2d.Tex())

	xf := sckorok.Transform.NewComp(entity)
	xf.SetPosition(f32.Vec2{200, 100})
}

func (*MainScene) Update(dt float32) {
}

func (*MainScene) OnExit() {
}

func main() {
	// Run game
	options := &sckorok.Options{
		Title:  "Simple Mesh Rendering",
		Width:  480,
		Height: 320,
	}
	sckorok.Run(options, &MainScene{})
}
