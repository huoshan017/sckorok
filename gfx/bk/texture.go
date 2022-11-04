package bk

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"sckorok/hid/gl"
	"sckorok/math/f32"
	"unsafe"
)

/**
处理纹理相关问题
*/

type Texture2D struct {
	Width, Height float32
	Id            uint32
}

func (t *Texture2D) Create(image image.Image) error {
	t.Width = float32(image.Bounds().Dx())
	t.Height = float32(image.Bounds().Dy())

	if id, err := newTexture(image); err != nil {
		return err
	} else {
		t.Id = id
	}
	return nil
}

func (t *Texture2D) Update(img image.Image, xoff, yoff int32, w, h int32) (err error) {
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		err = fmt.Errorf("unsupported stride")
		return
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.Id)
	gl.TexSubImage2D(gl.TEXTURE_2D,
		0,
		xoff,
		yoff,
		w,
		h,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&rgba.Pix[0]))
	return
}

func (t *Texture2D) Bind(stage int32) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(stage))
	gl.BindTexture(gl.TEXTURE_2D, t.Id)
}

func (t *Texture2D) Sub(x, y float32, w, h float32) *SubTex {
	subTex := &SubTex{Texture2D: t}
	subTex.Min = f32.Vec2{x, y}
	subTex.Max = f32.Vec2{x + w, y + h}
	return subTex
}

func (t *Texture2D) Destroy() {
	gl.DeleteTextures(1, &t.Id)
}

// TODO 提前转换图片格式
func newTexture(img image.Image) (uint32, error) {
	// 3. copy image
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	// 4. upload texture
	var texture uint32
	// 4.1 apply space
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// 4.2 params
	// 大小插值
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// 环绕方式
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// 4.3 upload
	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Dx()),
		int32(rgba.Rect.Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&rgba.Pix[0]))
	return texture, nil
}

///// 还需要抽象 SubTexture 的概念出来
type SubTex struct {
	*Texture2D

	// location -
	Min, Max f32.Vec2
}
