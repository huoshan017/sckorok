package main

import (
	korok "sckorok"
	"sckorok/asset"
	"sckorok/effect"
	"sckorok/game"
	"sckorok/math"
	"sckorok/math/f32"
)

type MainScene struct {
}

func (*MainScene) Load() {
	asset.Texture.Load("particle.png")
}

func (*MainScene) OnEnter(g *game.Game) {
	cfg := &effect.GravityConfig{
		Config: effect.Config{
			Max:      1024,
			Rate:     10,
			Duration: math.MaxFloat32,
			Life:     effect.Var{40.1, 0.4},
			Size:     effect.Range{effect.Var{10, 5}, effect.Var{20, 5}},
			X:        effect.Var{0, 0}, Y: effect.Var{0, 0},
			A: effect.Range{effect.Var{1, 0}, effect.Var{0, 0}},
		},
		Speed:   effect.Var{70, 10},
		Angel:   effect.Var{math.Radian(90), math.Radian(30)},
		Gravity: f32.Vec2{0, -10},
	}
	gravity := korok.Entity.New()
	gParticle := korok.ParticleSystem.NewComp(gravity)
	gParticle.SetSimulator(effect.NewGravitySimulator(cfg))
	gParticle.SetTexture(asset.Texture.Get("particle.png"))
	xf := korok.Transform.NewComp(gravity)
	xf.SetPosition(f32.Vec2{240, 160})
}

func (*MainScene) Update(dt float32) {

}

func (*MainScene) OnExit() {
}

func main() {
	options := &korok.Options{
		Title:  "ParticleSystem",
		Width:  480,
		Height: 320,
	}
	korok.Run(options, &MainScene{})
}
