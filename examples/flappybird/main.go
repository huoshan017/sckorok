package main

import (
	"sckorok"
	"sckorok/anim"
	"sckorok/asset"
	"sckorok/engi"
	"sckorok/game"
	"sckorok/gfx"
	"sckorok/gfx/font"
	"sckorok/gui"
	"sckorok/math/ease"
	"sckorok/math/f32"
)

type StartScene struct {
	title struct {
		gfx.Tex2D
		gui.Rect
	}
	start struct {
		btnNormal  gfx.Tex2D
		btnPressed gfx.Tex2D
		gui.Rect
	}
	bird, bg, ground engi.Entity
	mask             gfx.Color
}

func (sn *StartScene) Load() {
	asset.Texture.LoadAtlas("images/bird.png", "images/bird.json")
	asset.Font.LoadTrueType("font1", "fonts/Marker Felt.ttf", font.ASCII(32))
}

func (sn *StartScene) OnEnter(g *game.Game) {
	font, _ := asset.Font.Get("font1")
	gui.SetFont(font)

	at, _ := asset.Texture.Atlas("images/bird.png")
	bg, _ := at.GetByName("background.png")
	ground, _ := at.GetByName("ground.png")

	// setup gui
	// title
	tt, _ := at.GetByName("game_name.png")
	sn.title.Tex2D = tt
	sn.title.Rect = gui.Rect{
		X: (320 - 233) / 2,
		Y: 80,
		W: 233,
		H: 70,
	}

	// start button
	btn, _ := at.GetByName("start.png")
	sn.start.btnNormal = btn
	sn.start.btnPressed = btn
	sn.start.Rect = gui.Rect{
		X: (320 - 120) / 2,
		Y: 300,
		W: 120,
		H: 60,
	}

	// setup bg
	{
		entity := sckorok.Entity.New()
		spr := sckorok.Sprite.NewCompX(entity, bg)
		spr.SetSize(320, 480)
		xf := sckorok.Transform.NewComp(entity)
		xf.SetPosition(f32.Vec2{160, 240})
		sn.bg = entity
	}

	// setup ground {840 281}
	{
		entity := sckorok.Entity.New()
		spr := sckorok.Sprite.NewCompX(entity, ground)
		spr.SetSize(420, 140)
		spr.SetGravity(0, 1)
		spr.SetZOrder(1)
		xf := sckorok.Transform.NewComp(entity)
		xf.SetPosition(f32.Vec2{0, 100})
		sn.ground = entity
	}

	// flying animation
	bird1, _ := at.GetByName("bird1.png")
	bird2, _ := at.GetByName("bird2.png")
	bird3, _ := at.GetByName("bird3.png")

	frames := []gfx.Tex2D{bird1, bird2, bird3}
	g.AnimationSystem.SpriteEngine.NewAnimation("flying", frames, true)

	// setup bird
	bird := sckorok.Entity.New()
	spr := sckorok.Sprite.NewCompX(bird, bird1)
	spr.SetSize(48, 32)
	spr.SetZOrder(2)
	xf := sckorok.Transform.NewComp(bird)
	xf.SetPosition(f32.Vec2{160, 240})

	anim := sckorok.Flipbook.NewComp(bird)
	anim.SetRate(.1)
	anim.Play("flying")

	sn.bird = bird
}
func (sn *StartScene) Update(dt float32) {
	// draw title
	gui.Image(1, sn.title.Rect, sn.title.Tex2D, nil)

	// draw start button
	e := gui.ImageButton(2, sn.start.Rect, sn.start.btnNormal, sn.start.btnPressed, nil)
	if e.JustPressed() {
		// sn.LoadGame()
		sn.fadeOut()
	}
	// fade color
	if sn.mask.A > 0 {
		gui.ColorRect(gui.Rect{W: 320, H: 480}, sn.mask, 0)
	}
}
func (sn *StartScene) OnExit() {
}

func (sn *StartScene) fadeOut() {
	anim.OfColor(&sn.mask, gfx.Transparent, gfx.White).SetFunction(ease.InOutSine).SetDuration(1).OnComplete(func(reverse bool) {
		sn.loadGame()
	}).Forward()
}

func (sn *StartScene) loadGame() {
	gsn := &GameScene{}
	gsn.borrow(sn.bird, sn.bg, sn.ground)

	// load game scene
	sckorok.SceneMan.Load(gsn)
	sckorok.SceneMan.Push(gsn)
}

func main() {
	options := sckorok.Options{
		Title:  "Flappy Bird",
		Width:  320,
		Height: 480,
	}
	sckorok.Run(&options, &StartScene{})
}
