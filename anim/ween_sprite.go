package anim

import (
	"sckorok/anim/ween"
	"sckorok/engi"
	"sckorok/gfx"
)

// Convenient methods that uses to animate the Sprite Component.

// Tint the Entity to given color.
func Tint(e engi.Entity, from, to gfx.Color) *proxyAnimator {
	proxy := &proxyAnimator{Animator: tweenEngine.NewAnimator()}
	proxy.Animator.OnUpdate(func(reverse bool, f float32) {
		if spr := animationSystem.st.Comp(e); spr != nil {
			c := ween.ColorLerp(from, to, f)
			spr.SetColor(c)
		}
		if fn := proxy.update; fn != nil {
			fn(reverse, f)
		}
	})
	proxy.Animator.OnComplete(func(reverse bool) {
		proxy.Dispose()
		if fn := proxy.complete; fn != nil {
			fn(reverse)
		}
	})
	return proxy
}

//
//func Alpha(e engi.Entity, from, to float32) *proxyAnimator {
//	proxy := &proxyAnimator{Animator:tweenEngine.NewAnimator()}
//	proxy.Animator.OnUpdate(func(reverse bool, f float32) {
//		if spr := animationSystem.st.Comp(e); spr != nil {
//			c := spr.Color()
//			a := ween.F32Lerp(from, to, f)
//			// TODO
//		}
//		if fn := proxy.update; fn != nil {
//			fn(reverse, f)
//		}
//	})
//	proxy.Animator.OnComplete(func(reverse bool) {
//		proxy.Dispose()
//		if fn := proxy.complete; fn != nil {
//			fn(reverse)
//		}
//	})
//	return proxy
//}
