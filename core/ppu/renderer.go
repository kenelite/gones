package ppu

type Renderer struct {
	Framebuffer [240][256]byte // 每帧图像缓冲区
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

// 渲染一帧：可以先用随机颜色测试，再实现 tile 渲染
func (r *Renderer) RenderFrame(ppu *PPU) {
	for y := 0; y < 240; y++ {
		for x := 0; x < 256; x++ {
			r.Framebuffer[y][x] = byte((x + y + ppu.Frame) % 64)
		}
	}
}
