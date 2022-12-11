package main

import (
	"sckorok"
	"sckorok/anim/frame"
	"sckorok/asset"
	"sckorok/game"
	"sckorok/gfx"
	"sckorok/hid/input"
	"sckorok/math/f32"
)

type MainScene struct {
	heroAnim      *frame.FlipbookComp
	heroTransform *gfx.Transform
}

func (*MainScene) Load() {
	asset.Texture.LoadAtlasIndexed("hero.png", 52, 72, 4, 3)
}

func (m *MainScene) OnEnter(g *game.Game) {
	// get animation system...

	// input control
	input.RegisterButton("up", input.ArrowUp)
	input.RegisterButton("down", input.ArrowDown)
	input.RegisterButton("left", input.ArrowLeft)
	input.RegisterButton("right", input.ArrowRight)

	hero := sckorok.Entity.New()

	// SpriteComp
	sckorok.Sprite.NewComp(hero).SetSize(50, 50)
	sckorok.Transform.NewComp(hero).SetPosition(f32.Vec2{240, 160})

	fb := sckorok.Flipbook.NewComp(hero)
	fb.SetRate(.2)
	fb.SetLoop(true, frame.Restart)
	m.heroAnim = fb
	m.heroTransform = sckorok.Transform.Comp(hero)

	{
		at, _ := asset.Texture.Atlas("hero.png")
		frames := [12]gfx.Tex2D{}
		for i := 0; i < 12; i++ {
			frames[i], _ = at.GetByIndex(i)
		}
		g.SpriteEngine.NewAnimation("hero.down", frames[0:3], true)
		g.SpriteEngine.NewAnimation("hero.left", frames[3:6], true)
		g.SpriteEngine.NewAnimation("hero.right", frames[6:9], true)
		g.SpriteEngine.NewAnimation("hero.top", frames[9:12], true)
	}

	// default
	m.heroAnim.Play("hero.down")
}

func (m *MainScene) Update(dt float32) {
	speed := f32.Vec2{0, 0}

	// 根据上下左右，执行不同的帧动画
	if input.Button("up").JustPressed() {
		m.heroAnim.Play("hero.top")
	}
	if input.Button("down").JustPressed() {
		m.heroAnim.Play("hero.down")
	}
	if input.Button("left").JustPressed() {
		m.heroAnim.Play("hero.left")
	}
	if input.Button("right").JustPressed() {
		m.heroAnim.Play("hero.right")
	}

	scalar := float32(3)
	if input.Button("up").Down() {
		speed[1] = scalar
	}
	if input.Button("down").Down() {
		speed[1] = -scalar
	}
	if input.Button("left").Down() {
		speed[0] = -scalar
	}
	if input.Button("right").Down() {
		speed[0] = scalar
	}

	x := m.heroTransform.Position()[0] + speed[0]
	y := m.heroTransform.Position()[1] + speed[1]
	m.heroTransform.SetPosition(f32.Vec2{x, y})
}

func (*MainScene) OnExit() {
}

func main() {
	// Run game
	options := &sckorok.Options{
		Title:  "Hello, Korok Engine",
		Width:  480,
		Height: 320,
	}
	sckorok.Run(options, &MainScene{})
}
