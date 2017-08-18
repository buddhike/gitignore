package main

type Input struct {
	position int
	buffer   []rune
}

func (i *Input) current() (rune, bool) {
	if i.position == len(i.buffer) {
		return 0, true
	}

	return i.buffer[i.position], false
}

func (i *Input) advance() (rune, bool) {
	if i.position < len(i.buffer) {
		i.position++
	}

	return i.current()
}

func newInput(value string) Input {
	return Input{
		position: 0,
		buffer:   []rune(value),
	}
}

func (i *Input) last() (rune, bool) {
	if len(i.buffer) == 0 {
		return 0, false
	}
	return i.buffer[len(i.buffer)-1], true
}

func (i *Input) first() (rune, bool) {
	if len(i.buffer) == 0 {
		return 0, false
	}
	return i.buffer[0], true
}
