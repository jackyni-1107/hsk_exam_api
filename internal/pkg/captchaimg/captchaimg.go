package captchaimg

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	stdraw "image/draw"
	"image/png"
	"math/rand"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// QuestionToPNGBase64 将题目文案渲染为 PNG，返回标准 Base64（不含 data URL 前缀）。
// 使用随机底色、字色、整体缩放（字号感）及干扰线/噪点做轻量混淆。
func QuestionToPNGBase64(question string) (string, error) {
	if question == "" {
		question = "?"
	}
	face := basicfont.Face7x13
	padding := 10 + rand.Intn(9) // 10–18
	w := font.MeasureString(face, question)
	width := w.Ceil() + 2*padding
	if width < 96 {
		width = 96
	}
	height := face.Height + 2*padding
	if height < 40 {
		height = 40
	}

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	bg := randomLightBG()
	stdraw.Draw(rgba, rgba.Bounds(), image.NewUniform(bg), image.Point{}, stdraw.Src)

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(randomDarkFG()),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(padding), Y: fixed.I(padding + face.Ascent)},
	}
	d.DrawString(question)

	// 整体缩放，模拟不同字号（0.72–1.28）
	scale := 0.72 + rand.Float64()*0.56
	outW := maxInt(72, int(float64(width)*scale+0.5))
	outH := maxInt(36, int(float64(height)*scale+0.5))
	scaled := image.NewRGBA(image.Rect(0, 0, outW, outH))
	draw.CatmullRom.Scale(scaled, scaled.Bounds(), rgba, rgba.Bounds(), stdraw.Src, nil)

	addNoiseLines(scaled)
	addNoiseDots(scaled)

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func randomLightBG() color.RGBA {
	base := uint8(232 + rand.Intn(20))
	d := uint8(rand.Intn(10))
	return color.RGBA{
		R: base - d,
		G: base - uint8(rand.Intn(8)),
		B: base - uint8(rand.Intn(6)),
		A: 0xff,
	}
}

func randomDarkFG() color.RGBA {
	return color.RGBA{
		R: uint8(12 + rand.Intn(58)),
		G: uint8(12 + rand.Intn(58)),
		B: uint8(12 + rand.Intn(58)),
		A: 0xff,
	}
}

func addNoiseLines(img *image.RGBA) {
	b := img.Bounds()
	n := 3 + rand.Intn(4) // 3–6 条
	for i := 0; i < n; i++ {
		x0 := rand.Intn(b.Dx())
		y0 := rand.Intn(b.Dy())
		x1 := rand.Intn(b.Dx())
		y1 := rand.Intn(b.Dy())
		c := color.RGBA{
			R: uint8(100 + rand.Intn(120)),
			G: uint8(100 + rand.Intn(120)),
			B: uint8(100 + rand.Intn(120)),
			A: uint8(40 + rand.Intn(100)),
		}
		line(img, b.Min.X+x0, b.Min.Y+y0, b.Min.X+x1, b.Min.Y+y1, c)
	}
}

func addNoiseDots(img *image.RGBA) {
	b := img.Bounds()
	n := 35 + rand.Intn(46) // 35–80 点
	for i := 0; i < n; i++ {
		x := b.Min.X + rand.Intn(b.Dx())
		y := b.Min.Y + rand.Intn(b.Dy())
		c := color.RGBA{
			R: uint8(rand.Intn(200)),
			G: uint8(rand.Intn(200)),
			B: uint8(rand.Intn(200)),
			A: uint8(90 + rand.Intn(166)),
		}
		img.SetRGBA(x, y, c)
		if rand.Intn(2) == 0 && x+1 < b.Max.X {
			img.SetRGBA(x+1, y, c)
		}
	}
}

// Bresenham 线段，用于干扰线
func line(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	for {
		blendSet(img, x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 与背景 alpha 混合，避免线条完全盖住文字
func blendSet(img *image.RGBA, x, y int, c color.RGBA) {
	b := img.Bounds()
	if x < b.Min.X || x >= b.Max.X || y < b.Min.Y || y >= b.Max.Y {
		return
	}
	dst := img.RGBAAt(x, y)
	a := float64(c.A) / 255
	inv := 1 - a
	img.SetRGBA(x, y, color.RGBA{
		R: uint8(float64(dst.R)*inv + float64(c.R)*a),
		G: uint8(float64(dst.G)*inv + float64(c.G)*a),
		B: uint8(float64(dst.B)*inv + float64(c.B)*a),
		A: 0xff,
	})
}
