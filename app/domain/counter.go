package domain

type Counter struct {
	frame    int
	quadrant int
	page     int
	fps      int
}

func NewCounter() Counter {
	return Counter{
		frame:    0,
		quadrant: 0,
		page:     -1,
		fps:      0,
	}
}

func (c Counter) Frame() int {
	return c.frame
}

func (c Counter) Page() int {
	return c.page
}

func (c Counter) Quadrant() int {
	return c.quadrant
}

func (c *Counter) InitQuadrant() {
	c.quadrant = 0
}

func (c *Counter) IncQuadrant() {
	c.quadrant += 1
}

func (c *Counter) IncFrame() {
	c.frame += 1
}

func (c *Counter) IncPage() {
	c.page += 1
}

func (c *Counter) CanIgnore(fps int) bool {
	if c.fps != 0 && c.fps != fps {
		c.fps += 1
		return true
	} else {
		c.fps = 1
		return false
	}
}

func (c Counter) IsFirstQuadrantForFirstRow() bool {
	return c.quadrant == 0
}

func (c Counter) IsFirstQuadrantForAnyRow() bool {
	return (c.quadrant != 0) && (c.quadrant%3 == 0)
}

func (c Counter) IsStartPage(totalFrameInPage int) bool {
	return (c.frame == 0) || (c.frame%totalFrameInPage == 0)
}
