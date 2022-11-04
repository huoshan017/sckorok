package asset

import (
	"io"
	"sckorok/asset/res"
	"sckorok/effect"
	"sckorok/math"
	"sckorok/math/f32"

	"encoding/json"
	"log"
)

// 粒子系统配置文件管理
type ParticleConfigManager struct {
	repo map[string]refCount
}

func NewParticleConfigManager() *ParticleConfigManager {
	return &ParticleConfigManager{
		repo: make(map[string]refCount),
	}
}

func (pcm *ParticleConfigManager) Load(file string) {
	if rc, ok := pcm.repo[file]; ok {
		pcm.repo[file] = refCount{rc.ref, rc.cnt + 1}
	} else {
		ref, err := pcm.load(file)
		if err != nil {
			log.Println(err)
		} else {
			pcm.repo[file] = refCount{ref, 1}
		}
	}
}

func (pcm *ParticleConfigManager) Unload(file string) {
	if rc, ok := pcm.repo[file]; ok {
		if rc.cnt > 1 {
			pcm.repo[file] = refCount{rc.ref, rc.cnt - 1}
		} else {
			delete(pcm.repo, file)
		}
	}
}

func (pcm *ParticleConfigManager) Get(file string) (res interface{}, exist bool) {
	if rc, ok := pcm.repo[file]; ok {
		res = rc.ref
		exist = ok
	}
	return
}

func (pcm *ParticleConfigManager) load(file string) (ref interface{}, err error) {
	render, err := res.Open(file)
	if err != nil {
		return
	}
	defer render.Close()

	data, err := io.ReadAll(render)
	if err != nil {
		return
	}

	cfg := &psConfig{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return
	}

	var config *effect.Config
	if cfg.EmitterType == 0 {
		g := &effect.GravityConfig{}
		ref = g
		config = &g.Config
		g.Gravity = f32.Vec2{cfg.GravityX, cfg.GravityY}
		g.Speed = effect.Var{Base: cfg.Speed, Var: cfg.SpeedVar}

		ab, av := math.Radian(cfg.Angle), math.Radian(cfg.AngleVar)
		g.Angel = effect.Var{Base: ab, Var: av}

		g.RadialAcc = effect.Var{Base: cfg.RadialAccel, Var: cfg.RadialAccelVar}
		g.TangentialAcc = effect.Var{Base: cfg.TangentialAccel, Var: cfg.TangentialAccelVar}
		g.RotationIsDir = cfg.RotationIsDir
	} else {
		r := effect.RadiusConfig{}
		ref = r
		config = &r.Config
		r.Radius = effect.Range{
			Start: effect.Var{Base: cfg.StartRadius, Var: cfg.StartRadiusVar},
			End:   effect.Var{Base: cfg.EndRadius, Var: cfg.EndRadiusVar},
		}
		r.Angle = effect.Var{Base: cfg.Angle, Var: cfg.AngleVar}
		//r.AngleDelta = effect.Var{cfg.Angle}
	}

	// shared properties
	config.Max = cfg.MaxParticles
	config.Duration = cfg.Duration
	config.Life = effect.Var{Base: cfg.LifeSpan, Var: cfg.LifeSpanVar}
	config.X = effect.Var{Base: cfg.SourcePositionX, Var: cfg.SourcePositionVarX}
	config.Y = effect.Var{Base: cfg.SourcePositionY, Var: cfg.SourcePositionVarY}

	// size and spin
	config.Size = effect.Range{
		Start: effect.Var{Base: cfg.StartSize, Var: cfg.StartSizeVar},
		End:   effect.Var{Base: cfg.EndSize, Var: cfg.EndSizeVar},
	}
	config.Rot = effect.Range{
		Start: effect.Var{Base: cfg.StartSpin, Var: cfg.StartSpinVar},
		End:   effect.Var{Base: cfg.EndSpin, Var: cfg.EndSpinVar},
	}
	// color
	config.R = effect.Range{
		Start: effect.Var{Base: cfg.StartColorRed, Var: cfg.StartColorVarRed},
		End:   effect.Var{Base: cfg.EndColorRed, Var: cfg.EndColorVarRed},
	}
	config.G = effect.Range{
		Start: effect.Var{Base: cfg.StartColorGreen, Var: cfg.StartColorVarGreen},
		End:   effect.Var{Base: cfg.EndColorGreen, Var: cfg.EndColorVarGreen},
	}
	config.B = effect.Range{
		Start: effect.Var{Base: cfg.StartColorBlue, Var: cfg.StartColorVarBlue},
		End:   effect.Var{Base: cfg.EndColorBlue, Var: cfg.EndColorVarBlue},
	}
	config.A = effect.Range{
		Start: effect.Var{Base: cfg.StartColorAlpha, Var: cfg.StartColorVarAlpha},
		End:   effect.Var{Base: cfg.EndColorAlpha, Var: cfg.EndColorVarAlpha},
	}
	return
}

type psConfig struct {
	ConfigName string `json:"configName"`

	MaxParticles int     `json:"maxParticles"`
	Angle        float32 `json:"angle"`
	AngleVar     float32 `json:"angleVariance"`
	Duration     float32 `json:"duration"`

	// blend-func - not support, now

	// color
	StartColorRed   float32 `json:"startColorRed"`
	StartColorGreen float32 `json:"startColorGreen"`
	StartColorBlue  float32 `json:"startColorBlue"`
	StartColorAlpha float32 `json:"startColorAlpha"`

	StartColorVarRed   float32 `json:"startColorVarianceRed"`
	StartColorVarGreen float32 `json:"startColorVarianceGreen"`
	StartColorVarBlue  float32 `json:"startColorVarianceBlue"`
	StartColorVarAlpha float32 `json:"startColorVarianceAlpha"`

	EndColorRed   float32 `json:"finishColorRed"`
	EndColorGreen float32 `json:"finishColorGreen"`
	EndColorBlue  float32 `json:"finishColorBlue"`
	EndColorAlpha float32 `json:"finishColorAlpha"`

	EndColorVarRed   float32 `json:"finishColorVarianceRed"`
	EndColorVarGreen float32 `json:"finishColorVarianceGreen"`
	EndColorVarBlue  float32 `json:"finishColorVarianceBlue"`
	EndColorVarAlpha float32 `json:"finishColorVarianceAlpha"`

	// size
	StartSize    float32 `json:"startParticleSize"`
	StartSizeVar float32 `json:"startParticleSizeVariance"`
	EndSize      float32 `json:"finishParticleSize"`
	EndSizeVar   float32 `json:"finishParticleSizeVariance"`

	// Position
	SourcePositionX float32 `json:"sourcePositionx"`
	SourcePositionY float32 `json:"sourcePositiony"`

	SourcePositionVarX float32 `json:"sourcePositionVariancex"`
	SourcePositionVarY float32 `json:"sourcePositionVariancey"`

	// Spinning
	StartSpin    float32 `json:"rotationStart"`
	StartSpinVar float32 `json:"rotationStartVariance"`
	EndSpin      float32 `json:"rotationEnd"`
	EndSpinVar   float32 `json:"rotationEndVariance"`

	// life and emission rate
	LifeSpan    float32 `json:"particleLifespan"`
	LifeSpanVar float32 `json:"particleLifespanVariance"`

	// mode
	EmitterType int `json:"emitterType"`

	///////// ModeA

	// gravity
	GravityX float32 `json:"gravityx"`
	GravityY float32 `json:"gravityy"`

	// speed
	Speed    float32 `json:"speed"`
	SpeedVar float32 `json:"speedVariance"`

	// radial acceleration
	RadialAccel    float32 `json:"radialAcceleration"`
	RadialAccelVar float32 `json:"radialAccelVariance"`

	// tangential acceleration
	TangentialAccel    float32 `json:"tangentialAcceleration"`
	TangentialAccelVar float32 `json:"tangentialAccelVariance"`

	// rotation is dir
	RotationIsDir bool `json:"rotationIsDir"`

	////////// ModeB

	// radius
	StartRadius    float32 `json:"maxRadius"`
	StartRadiusVar float32 `json:"maxRadiusVariance"`

	EndRadius    float32 `json:"minRadius"`
	EndRadiusVar float32 `json:"minRadiusVariance"`

	// rotate
	RotatePerSecond    float32 `json:"rotatePerSecond"`
	RotatePerSecondVar float32 `json:"rotatePerSecondVariance"`
}
