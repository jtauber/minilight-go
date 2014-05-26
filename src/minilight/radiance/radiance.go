package radiance

import (
	"image"
	"image/color"
	"math"
)

const R_LUMINANCE = 0.2126
const G_LUMINANCE = 0.7152
const B_LUMINANCE = 0.0722

const DISPLAY_LUMINANCE_MAX = 200.0

const GAMMA_ENCODE = 0.45

type Radiance struct {
	R, G, B float64
}

type RadianceImage struct {
	Width, Height int
	pix           []float64
}

func NewRadianceImage(width, height int) *RadianceImage {
	buf := make([]float64, 3*width*height)
	return &RadianceImage{width, height, buf}
}

func (p *RadianceImage) pixOffset(x, y int) int {
	return y*p.Width*3 + x*3
}

func (p *RadianceImage) Add(x int, y int, radiance Radiance) {
	i := p.pixOffset(x, y)
	p.pix[i+0] += radiance.R
	p.pix[i+1] += radiance.G
	p.pix[i+2] += radiance.B
}

func (p *RadianceImage) Get(x int, y int) Radiance {
	i := p.pixOffset(x, y)
	return Radiance{p.pix[i + 0], p.pix[i + 1], p.pix[i + 2]}
}

func (p *RadianceImage) calculateScaleFactor(iterations int) float64 {
    // calculate the linear tone-mapping scalefactor for this image assuming
    // the given number of iterations.

    // calculate the log-mean luminance of the image

	var sum_of_logs float64

	for x := 0; x < p.Width; x++ {
		for y := 0; y < p.Height; y++ {
            lum := p.pix[p.pixOffset(x, y) + 0] * R_LUMINANCE
            lum += p.pix[p.pixOffset(x, y) + 1] * G_LUMINANCE
            lum += p.pix[p.pixOffset(x, y) + 2] * B_LUMINANCE
            lum /= float64(iterations)

            sum_of_logs += math.Log10(math.Max(lum, 0.0001))
		}
	}

    log_mean_luminance := math.Pow(10.0, (sum_of_logs / float64(p.Height * p.Width)))

    // calculate the scalefactor for linear tone-mapping
    // formula from Ward "A Contrast-Based Scalefactor for Luminance Display"

	scalefactor_numerator := 1.219 + math.Pow(DISPLAY_LUMINANCE_MAX * 0.25, 0.4)

    scalefactor := (math.Pow((scalefactor_numerator / (1.219 + math.Pow(log_mean_luminance, 0.4))), 2.5)) / DISPLAY_LUMINANCE_MAX

    return scalefactor
}

func (p *RadianceImage) ToRGBA64(iterations int) *image.RGBA64 {

	scalefactor := p.calculateScaleFactor(iterations)

	m := image.NewRGBA64(image.Rect(0, 0, p.Width, p.Height))
	for x := 0; x < p.Width; x++ {
		for y := 0; y < p.Height; y++ {
			radiance := p.Get(x, y)
			r := uint16(65535 * math.Pow(math.Max(radiance.R * scalefactor / float64(iterations), 0), GAMMA_ENCODE))
			g := uint16(65535 * math.Pow(math.Max(radiance.G * scalefactor / float64(iterations), 0), GAMMA_ENCODE))
			b := uint16(65535 * math.Pow(math.Max(radiance.B * scalefactor / float64(iterations), 0), GAMMA_ENCODE))
			a := uint16(65535)
			m.Set(x, y, color.RGBA64{r, g, b, a})
		}
	}
	return m
}
