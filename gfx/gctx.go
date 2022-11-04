package gfx

import (
	"sckorok/gfx/bk"
	"sort"
	"unsafe"
)

// graphics context
// a wrapper for bk-api

const SharedIndexBufferSize uint16 = 0xFFFF

func Init(pixelRatio float32) {
	bk.Init()
	bk.Reset(480, 320, pixelRatio)

	// Enable debug text
	bk.SetDebug(bk.DebugResMan | bk.DebugQueue)
}

func Flush() (num int) {
	num = bk.Flush()
	Context.Step()
	return
}

func Destroy() {
	bk.Destroy()
	Context.Destroy()
}

// 目前各个 RenderFeature 都是自己管理 VBO/IBO，但是对于一些系统，比如
// Batch/ParticleSystem(2D中的大部分元素)，都是可以复用VBO的，顶点数据
// 需要每帧动态生成，如此可以把这些需要动态申请的Buffer在此管理起来，对应的
// CPU 数据可以在 StackAllocator 上申请，一帧之后就自动释放。
type tempBuffer struct {
	vb     *bk.VertexBuffer
	size   int
	stride int
	id     uint16
	use    uint16
}

type context struct {
	Stack StackAllocator
	temps []tempBuffer

	shared struct {
		id      uint16
		padding uint16
		index   []uint16
		size    int
	}
}

// 一帧之后自动释放
func (ctx *context) TempVertexBuffer(reqSize, stride int) (id uint16, size int, vb *bk.VertexBuffer) {
	var (
		buffer tempBuffer
		found  bool
	)
	for i, tb := range ctx.temps {
		if tb.use == 0 && tb.stride == stride && tb.size >= reqSize {
			ctx.temps[i].use = 1
			buffer = tb
			found = true
			break
		}
	}
	if !found {
		buffer = ctx.newVertexBuffer(reqSize, stride)
		ctx.temps = append(ctx.temps, buffer)
	}
	id, size, vb = buffer.id, buffer.size, buffer.vb
	return
}

func (ctx *context) newVertexBuffer(vertexSize, stride int) tempBuffer {
	{
		vertexSize--
		vertexSize |= vertexSize >> 1
		vertexSize |= vertexSize >> 2
		vertexSize |= vertexSize >> 3
		vertexSize |= vertexSize >> 8
		vertexSize |= vertexSize >> 16
		vertexSize++
	}
	tb := tempBuffer{size: vertexSize, stride: stride, use: 1}
	if id, vb := bk.R.AllocVertexBuffer(bk.Memory{nil, uint32(vertexSize * stride)}, uint16(stride)); id != bk.InvalidId {
		tb.id = id
		tb.vb = vb
	}
	return tb
}

func (ctx *context) release() {
	for i := range ctx.temps {
		ctx.temps[i].use = 0
	}
	sort.Slice(ctx.temps, func(i, j int) bool {
		return ctx.temps[i].size < ctx.temps[j].size
	})
}

func (ctx *context) Step() {
	ctx.release()
	ctx.Stack.release()
}

func (ctx *context) Destroy() {
	// destroy any release!
}

// 64kb, format={3, 0, 1, 3, 1, 2}
func (ctx *context) SharedIndexBuffer() (id uint16, size int) {
	if ctx.shared.index == nil {
		ctx.initIndexBuffer()
	}
	return ctx.shared.id, ctx.shared.size
}

func (ctx *context) initIndexBuffer() {
	indexSize := int(SharedIndexBufferSize)
	ctx.shared.index = make([]uint16, indexSize)
	size := int(indexSize)
	iFormat := [6]uint16{3, 0, 1, 3, 1, 2}
	for i := 0; i < size; i += 6 {
		copy(ctx.shared.index[i:], iFormat[:])
		iFormat[0] += 4
		iFormat[1] += 4
		iFormat[2] += 4
		iFormat[3] += 4
		iFormat[4] += 4
		iFormat[5] += 4
	}
	if id, _ := bk.R.AllocIndexBuffer(bk.Memory{unsafe.Pointer(&ctx.shared.index[0]), uint32(size) * 2}); id != bk.InvalidId {
		ctx.shared.id = id
		ctx.shared.size = size
	}
}

// global shared
var Context *context

// Stack capacity ~ 100K
func init() {
	Context = &context{}
	Context.Stack.initialize(100 * 1024)
}
