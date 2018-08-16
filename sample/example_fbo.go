package main

import (
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"nanovgo"
)

func main() {
	var winWidth, winHeight int
	var fbWidth, fbHeight int
	var pixelRatio float32
	var fb *nanovgo.FrameBuffer
	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// demo MSAA
	glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(1000, 600, "NanoVGo", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	ctx, err := nanovgo.NewContext(nanovgo.AntiAlias | nanovgo.StencilStrokes | nanovgo.Debug)
	defer ctx.Delete()

	if err != nil {
		panic(err)
	}

	fbWidth, fbHeight = window.GetFramebufferSize()
	winWidth, winHeight = window.GetSize()

	pixelRatio = float32(fbWidth) / float32(winWidth)

	fb = ctx.CreateFramebuffer(fbWidth, fbHeight, 0)
	if fb == nil {
		panic("framebufferobject should be somewhere around there")
		return
	}

	glfw.SwapInterval(0)

	for !window.ShouldClose() {
		renderPattern(ctx, fb, pixelRatio)

		gl.Viewport(0, 0, fbWidth, fbHeight)
		gl.ClearColor(0, 0, 0, 0)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

		fbWidth, fbHeight = window.GetFramebufferSize()
		winWidth, winHeight = window.GetSize()

		pixelRatio = float32(fbWidth) / float32(winWidth)

		ctx.BeginFrame(winWidth, winHeight, pixelRatio)

		if fb != nil {
			ctx.Save()

			ctx.BeginPath()
			ctx.RoundedRect(300, 200, 200,200,10)
			ctx.Circle(200, 400, 50)
			ctx.SetFillColor(nanovgo.RGBA(255, 0, 0, 255))
			ctx.Fill()

			img := nanovgo.ImagePattern(0, 0, float32(fbWidth), float32(fbHeight), 0, fb.Image(), 1.0)
			ctx.BeginPath()
			ctx.Circle(200, 200, 150)
			ctx.SetFillPaint(img)
			ctx.Fill()
			ctx.Stroke()
			ctx.SetStrokeColor(nanovgo.RGBA(0, 255, 0, 255))
			ctx.SetStrokeWidth(3)
			ctx.Stroke()

			ctx.Restore()
		}

		ctx.EndFrame()

		gl.Enable(gl.DEPTH_TEST)
		window.SwapBuffers()
		glfw.PollEvents()
	}

	nanovgo.NvgluDeleteFramebuffer(fb)
}

func renderPattern(ctx *nanovgo.Context, fb *nanovgo.FrameBuffer, pxRatio float32) {
	var winWidth, winHeight int
	var fboWidth, fboHeight int
	//s := 20.0

	if fb == nil {
		return
	}

	fboWidth, fboHeight, _ = ctx.ImageSize(fb.Image())
	winWidth = int(float32(fboWidth) / pxRatio)
	winHeight = int(float32(fboHeight) / pxRatio)

	// Draw some stuff to an FBO as a test
	nanovgo.NvgluBindFramebuffer(fb)
	gl.Viewport(0, 0, fboWidth, fboHeight)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

	ctx.BeginFrame(winWidth, winHeight, pxRatio)

	for i := 0; i < 20; i++ {
		ctx.BeginPath()
		ctx.Rect(float32(10+i*30), 10, 5, float32(winHeight-20))
		ctx.Rotate(nanovgo.DegToRad(1))
		ctx.SetFillColor(nanovgo.HSLA(float32(i/19.0), 0.5, 0.5, 255))
		ctx.Fill()
	}

	ctx.EndFrame()

	nanovgo.NvgluBindFramebuffer(nil)
}
