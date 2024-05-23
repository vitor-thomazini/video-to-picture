package domain

const (
	MAX_VALUE_FRAMES_COUNTER = 18 //12
)

type FrameController struct {
	framesCounter     int
	startInIndexFrame int
	endInIndexFrame   int
	latestText        string
}

func NewFrameController(startInIndexFrame int, indexFrameDelay int) FrameController {
	return FrameController{
		framesCounter:     0,
		startInIndexFrame: startInIndexFrame,
		endInIndexFrame:   startInIndexFrame + indexFrameDelay,
	}
}

func (c *FrameController) Counter() int {
	return c.framesCounter
}

func (c *FrameController) Wait() bool {
	wait := !(c.framesCounter%MAX_VALUE_FRAMES_COUNTER == 0 && c.startInIndexFrame < c.framesCounter)
	c.framesCounter += 1
	return wait
}

func (c *FrameController) InitLatestText(text string) {
	if c.latestText == "" {
		c.latestText = text
	}
}

func (c *FrameController) UpdateLatestText(text string) {
	if c.latestText != "none" {
		c.latestText = text
	}
}

func (c FrameController) LatestText() string {
	return c.latestText
}

func (c FrameController) EndInIndexFrame() int {
	return c.endInIndexFrame
}
