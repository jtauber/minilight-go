package radiance

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
