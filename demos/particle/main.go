package main

import (
	"sckorok"
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
			Life:     effect.Var{Base: 40.1, Var: 0.4},
			Size:     effect.Range{Start: effect.Var{Base: 10, Var: 5}, End: effect.Var{Base: 20, Var: 5}},
			X:        effect.Var{Base: 0, Var: 0}, Y: effect.Var{Base: 0, Var: 0},
			A: effect.Range{Start: effect.Var{Base: 1, Var: 0}, End: effect.Var{Base: 0, Var: 0}},
		},
		Speed:   effect.Var{Base: 70, Var: 10},
		Angel:   effect.Var{Base: math.Radian(90), Var: math.Radian(30)},
		Gravity: f32.Vec2{0, -10},
	}
	gravity := sckorok.Entity.New()
	gParticle := sckorok.ParticleSystem.NewComp(gravity)
	gParticle.SetSimulator(effect.NewGravitySimulator(cfg))
	gParticle.SetTexture(asset.Texture.Get("particle.png"))
	xf := sckorok.Transform.NewComp(gravity)
	xf.SetPosition(f32.Vec2{240, 160})
}

func (*MainScene) Update(dt float32) {

}

func (*MainScene) OnExit() {
}

func main() {
	options := &sckorok.Options{
		Title:  "ParticleSystem",
		Width:  480,
		Height: 320,
	}
	sckorok.Run(options, &MainScene{})
}
