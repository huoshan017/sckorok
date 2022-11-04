package bk

import (
	"log"
	"sckorok/math/f32"
	"sort"
	"unsafe"
)

type SortMode int

const (
	Sequential SortMode = iota
	Ascending
	Descending
)

type Rect struct {
	x, y uint16
	w, h uint16
}

func (r *Rect) clear() {
	r.x, r.y = 0, 0
	r.w, r.h = 0, 0
}

func (r *Rect) isZero() bool {
	u64 := (*uint64)(unsafe.Pointer(r))
	return *u64 == 0
}

type Stream struct {
	vertexBuffer uint16
	vertexFormat uint16 // Offset | Stride， not used now!!

	firstVertex uint16
	numVertex   uint16
}

type RenderDraw struct {
	indexBuffer   uint16
	vertexBuffers [2]Stream
	textures      [2]uint16

	// index params
	firstIndex, num uint16

	// uniform range
	uniformBegin uint16
	uniformEnd   uint16

	// stencil and scissor
	stencil uint32
	scissor uint16

	// required renderer state
	state uint64
}

func (rd *RenderDraw) reset() {
	rd.indexBuffer = 0
	rd.firstIndex, rd.num = 0, 0
	rd.scissor = 0
}

// ~ 8000 draw call
const MAX_QUEUE_SIZE = 8 << 10

type RenderQueue struct {
	SortMode
	// render list
	sortKey    [MAX_QUEUE_SIZE]uint64
	sortValues [MAX_QUEUE_SIZE]uint16

	drawCallList [MAX_QUEUE_SIZE]RenderDraw
	drawCallNum  uint16

	sk SortKey

	// per-drawCall state cache
	drawCall RenderDraw

	uniformBegin uint16
	uniformEnd   uint16

	// per-frame state
	viewports [4]Rect
	scissors  [4]Rect
	clears    [4]struct {
		index   [8]uint8
		rgba    uint32
		depth   float32
		stencil uint8
		flags   uint16
	}

	// per-frame data flow
	rm *ResManager
	ub *UniformBuffer

	// render context
	ctx *RenderContext
}

func NewRenderQueue(m *ResManager) *RenderQueue {
	ub := NewUniformBuffer()
	rc := NewRenderContext(m, ub)
	return &RenderQueue{
		ctx: rc,
		rm:  m,
		ub:  ub,
	}
}

func (rq *RenderQueue) Init() {
	rq.ctx.Init()
}

// reset frame-buffer size
func (rq *RenderQueue) Reset(w, h uint16, pr float32) {
	rq.ctx.wRect.w = w
	rq.ctx.wRect.h = h
	rq.ctx.pixelRatio = pr
}

func (rq *RenderQueue) Destroy() {
	//
}

func (rq *RenderQueue) SetState(state uint64, rgba uint32) {
	rq.drawCall.state = state
}

func (rq *RenderQueue) SetIndexBuffer(id uint16, firstIndex, num uint16) {
	rq.drawCall.indexBuffer = id & IdMask
	rq.drawCall.firstIndex = firstIndex
	rq.drawCall.num = num
}

func (rq *RenderQueue) SetVertexBuffer(stream uint8, id uint16, firstVertex, numVertex uint16) {
	if stream < 0 || stream >= 2 {
		log.Printf("Not support stream location: %d", stream)
		return
	}

	vbStream := &rq.drawCall.vertexBuffers[stream]
	vbStream.vertexBuffer = id & IdMask
	vbStream.vertexFormat = InvalidId
	vbStream.firstVertex = firstVertex
	vbStream.numVertex = numVertex
}

func (rq *RenderQueue) SetTexture(stage uint8, samplerId uint16, texId uint16, flags uint32) {
	if stage < 0 || stage >= 2 {
		log.Printf("Not suppor texture location: %d", stage)
		return
	}

	rq.drawCall.textures[stage] = texId & IdMask
}

// 复制简单数据的时候（比如：Sampler），采用赋值的方式可能更快 TODO
func (rq *RenderQueue) SetUniform(id uint16, ptr unsafe.Pointer) {
	if ok, um := rq.rm.Uniform(id); ok {
		opCode := Uniform_encode(um.Type, um.Slot, um.Size, um.Count)
		rq.ub.WriteUInt32(opCode)
		rq.ub.Copy(ptr, uint32(um.Size)*uint32(um.Count))
	}
}

// Transform 是 uniform 之一，2D 世界可以省略
func (rq *RenderQueue) SetTransform(mtx *f32.Mat4) {
	// TODO impl
}

func (rq *RenderQueue) SetStencil(stencil uint32) {
	rq.drawCall.stencil = stencil
}

func (rq *RenderQueue) SetScissor(x, y, width, height uint16) (id uint16) {
	id = rq.ctx.AddClipRect(x, y, width, height)
	rq.drawCall.scissor = id
	return id
}

func (rq *RenderQueue) SetScissorCached(id uint16) {
	rq.drawCall.scissor = id
}

/// View Related Setting
func (rq *RenderQueue) SetViewScissor(id uint8, x, y, with, height uint16) {
	if id < 0 || id >= 4 {
		log.Printf("Not support view id: %d", id)
		return
	}
	rq.scissors[id] = Rect{x, y, with, height}
}

func (rq *RenderQueue) SetViewPort(id uint8, x, y, width, height uint16) {
	if id < 0 || id >= 4 {
		log.Printf("Not support view id: %d", id)
		return
	}
	rq.viewports[id] = Rect{x, y, width, height}
}

func (rq *RenderQueue) SetViewClear(id uint8, flags uint16, rgba uint32, depth float32, stencil uint8) {
	if id < 0 || id >= 4 {
		log.Printf("Not support view id: %d", id)
		return
	}
	clear := &rq.clears[id]
	clear.flags = flags
	clear.rgba = rgba
	clear.depth = depth
	clear.stencil = stencil
}

func (rq *RenderQueue) SetViewTransform(id uint8, view, proj *f32.Mat4, flags uint8) {

}

// conversion: depth range [-int16, int16]
func (rq *RenderQueue) Submit(id uint8, program uint16, depth int32) uint32 {
	// uniform range
	rq.uniformEnd = uint16(rq.ub.GetPos())

	// encode sort-key
	sk := &rq.sk
	sk.Layer = uint16(id)
	sk.Order = uint16(depth + 0xFFFF>>1)

	sk.Shader = program & IdMask // trip type
	sk.Blend = 0
	sk.Texture = rq.drawCall.textures[0]

	rq.sortKey[rq.drawCallNum] = rq.sk.Encode()
	rq.sortValues[rq.drawCallNum] = rq.drawCallNum

	// copy data
	rq.drawCall.uniformBegin = rq.uniformBegin
	rq.drawCall.uniformEnd = rq.uniformEnd

	rq.drawCallList[rq.drawCallNum] = rq.drawCall
	rq.drawCallNum++

	// reset state
	rq.drawCall.reset()
	rq.uniformBegin = uint16(rq.ub.GetPos())

	// return frame Num
	return 0
}

/// 执行最终的绘制
func (rq *RenderQueue) Flush() int {
	num := rq.drawCallNum

	var (
		sortKeys = rq.sortKey[:num]
		sortVals = rq.sortValues[:num]
		drawList = rq.drawCallList[:num]
	)

	// Sort by SortKey
	switch rq.SortMode {
	case Ascending:
		sort.Stable(ByKeyAscending{sortKeys, sortVals})
	case Descending:
		sort.Stable(ByKeyDescending{sortKeys, sortVals})
	}

	// Draw respect to sorted values
	rq.ctx.Draw(sortKeys, sortVals, drawList)

	// Clear counter
	rq.drawCallNum = 0
	rq.uniformBegin = 0
	rq.uniformEnd = 0
	rq.ub.Reset()
	rq.ctx.Reset()

	return int(num)
}

// For 2D games, batch-system will reduce draw-call obviously,
// A simple default sort method will be OK.

// sort draw-call based on SortKey
type ByKeyAscending struct {
	k []uint64
	v []uint16
}

func (a ByKeyAscending) Len() int {
	return len(a.k)
}

func (a ByKeyAscending) Swap(i, j int) {
	a.k[i], a.k[j] = a.k[j], a.k[i]
	a.v[i], a.v[j] = a.v[j], a.v[i]
}

func (a ByKeyAscending) Less(i, j int) bool {
	return a.k[i] < a.k[j]
}

type ByKeyDescending struct {
	k []uint64
	v []uint16
}

func (a ByKeyDescending) Len() int {
	return len(a.k)
}

func (a ByKeyDescending) Swap(i, j int) {
	a.k[i], a.k[j] = a.k[j], a.k[i]
	a.v[i], a.v[j] = a.v[j], a.v[i]
}

func (a ByKeyDescending) Less(i, j int) bool {
	return a.k[i] > a.k[j]
}
