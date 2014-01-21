package radiance

type Radiance struct {
	R, G, B float64
}

type RadianceImage struct {
	Width, Height int
	Pix           []float64
}

func (p *RadianceImage) PixOffset(x, y int) int {
	return y*p.Width*3 + x*3
}

func (p *RadianceImage) Add(x int, y int, radiance Radiance) {
	i := p.PixOffset(x, y)
	p.Pix[i+0] += radiance.R
	p.Pix[i+1] += radiance.G
	p.Pix[i+2] += radiance.B
}

func NewRadianceImage(width, height int) *RadianceImage {
	buf := make([]float64, 3*width*height)
	return &RadianceImage{width, height, buf}
}
