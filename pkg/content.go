package pkg

type Content struct {
	Buffer [][]rune
}

func NewContent() *Content {
	buffer := make([][]rune, 1)
	buffer[0] = make([]rune, 0)
	return &Content{
		Buffer: buffer,
	}
}

func (c *Content) LineLength(lineNumber int) int {
	return len(c.Buffer[lineNumber])
}

func (c *Content) TotalLines() int {
	return len(c.Buffer)
}

func (c *Content) ToString() string {
	fileContentV2 := ""
	for _, line := range c.Buffer {
		fileContentV2 += string(line) + "\n"
	}
	return fileContentV2
}
