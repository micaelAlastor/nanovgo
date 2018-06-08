package nanovgo

import "github.com/goxjs/gl"

type GLint = int32
type GLuint = uint32

var defaultFBO = gl.Framebuffer{0}

type NVGLUFrameBuffer struct {
	ctx *glContext
	fbo gl.Framebuffer
	rbo gl.Renderbuffer
	texture *glTexture
	image int
}

func GetBoundRenderbuffer() gl.Renderbuffer {
	var b int
	b = gl.GetInteger(gl.RENDERBUFFER_BINDING)
	return gl.Renderbuffer{Value: uint32(b)}
}

func (p *glParams) RenderCreateFramebuffer(w, h int, flags ImageFlags) *NVGLUFrameBuffer {

	var defaultFBO gl.Framebuffer
	var defaultRBO gl.Renderbuffer

	defaultFBO = gl.GetBoundFramebuffer()
	defaultRBO = GetBoundRenderbuffer()

	fb := new(NVGLUFrameBuffer)
	fb.image = p.renderCreateTexture(nvgTextureRGBA, w, h, flags, nil)
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

	gl.BindFramebuffer(gl.FRAMEBUFFER, defaultFBO)
	gl.BindRenderbuffer(gl.RENDERBUFFER, defaultRBO)

	return fb
}

func NvgluBindFramebuffer(fb *NVGLUFrameBuffer){
	if !defaultFBO.Valid() {
		defaultFBO = gl.GetBoundFramebuffer()
	}
	if fb != nil {
		gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo)
	} else {
		gl.BindFramebuffer(gl.FRAMEBUFFER, defaultFBO)
	}

}

func NvgluDeleteFramebuffer(fb *NVGLUFrameBuffer) {
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

