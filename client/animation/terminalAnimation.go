package animation

type TerminalAnimation struct {
	NonBlockingAnimation
}

func (ta *TerminalAnimation) Next() Animation {
	inc := ta.cur + ta.animationSpeed.IncFrames()
	if inc >= len(ta.Sprites) {
		return nil
	}
	ta.cur = inc % len(ta.Sprites)
	return ta
}

func (ta *TerminalAnimation) ChangeAnimation(animation Animation) Animation {
	return ta
}

func (ta *TerminalAnimation) Copy() Animation{
	cpy := *ta
	return &cpy
}