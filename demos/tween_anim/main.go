package main

import (
	korok "sckorok"
	"sckorok/anim"
	"sckorok/anim/ween"
	"sckorok/asset"
	"sckorok/engi"
	"sckorok/game"
	"sckorok/math/ease"
	"sckorok/math/f32"
)

type MainScene struct {
	hero engi.Entity
}

func (*MainScene) Load() {
	asset.Texture.Load("face.png")
}

func (m *MainScene) OnEnter(g *game.Game) {
	// texture
	tex := asset.Texture.Get("face.png")

	// ease functions
	funcs := []ease.Function{
		ease.Linear,
		ease.OutCirc,
		ease.OutBounce,
		ease.OutElastic,
		ease.OutBack,
		ease.OutCubic,
	}

	for i := range funcs {
		entity := korok.Entity.New()
		korok.Sprite.NewCompX(entity, tex).SetSize(30, 30)
		korok.Transform.NewComp(entity).SetPosition(f32.Vec2{0, 50 + 30*float32(i)})
		anim.MoveX(entity, 10, 240).SetFunction(funcs[i]).SetRepeat(ween.RepeatInfinite, ween.Restart).SetDuration(2).Forward()
	}
}

func (m *MainScene) Update(dt float32) {
}

func (*MainScene) OnExit() {
}

func main() {
	// Run game
	options := &korok.Options{
		Title:  "Hello, Korok Engine",
		Width:  480,
		Height: 320,
	}
	korok.Run(options, &MainScene{})
}
