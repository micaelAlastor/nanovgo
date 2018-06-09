package nanovgo

import (
	"github.com/goxjs/gl"
	"fmt"
)

type GLint = int32
type GLuint = uint32

var defaultFBO = gl.Framebuffer{0}

type FrameBuffer struct {
	ctx *glContext
	fbo gl.Framebuffer
	rbo gl.Renderbuffer
	texture *glTexture
	image int
}

func (fb *FrameBuffer) Image() int {
	return fb.image
}

func GetBoundRenderbuffer() gl.Renderbuffer {
	var b int
	b = gl.GetInteger(gl.RENDERBUFFER_BINDING)
	return gl.Renderbuffer{Value: uint32(b)}
}

func (p *glParams) renderCreateFramebuffer(w, h int, flags ImageFlags) *FrameBuffer {

	var defaultFBO gl.Framebuffer
	var defaultRBO gl.Renderbuffer

	defaultFBO = gl.GetBoundFramebuffer()
	defaultRBO = GetBoundRenderbuffer()

	fb := new(FrameBuffer)
	fb.image = p.renderCreateTexture(nvgTextureRGBA, w, h, flags|ImagePreMultiplied, nil)
	fb.texture = p.context.findTexture(fb.image)
	fb.ctx = p.context

	// frame buffer object
	fb.fbo = gl.CreateFramebuffer()
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo)

	// render buffer object
	fb.rbo = gl.CreateRenderbuffer()
	gl.BindRenderbuffer(gl.RENDERBUFFER, fb.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.STENCIL_INDEX8, w, h)

	// combine all
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fb.texture.tex, 0)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.STENCIL_ATTACHMENT, gl.RENDERBUFFER, fb.rbo)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE{
		fmt.Println("YOBA")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, defaultFBO)
	gl.BindRenderbuffer(gl.RENDERBUFFER, defaultRBO)


	return fb
}

func NvgluBindFramebuffer(fb *FrameBuffer){
	if fb != nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo)
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, defaultFBO)
	}

}

func NvgluDeleteFramebuffer(fb *FrameBuffer) {
	if fb == nil {
		return
	}
	if fb.fbo.Valid(){
		gl.DeleteFramebuffer(fb.fbo)
	}
	if fb.rbo.Valid(){
		gl.DeleteRenderbuffer(fb.rbo)
	}
	if fb.image >= 0 {
		fb.ctx.deleteTexture(fb.image)
	}
	fb.ctx = nil
	fb.fbo.Value = 0
	fb.rbo.Value = 0
	fb.texture = nil
	fb.image = -1
}

