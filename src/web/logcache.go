package web

import "strings"

type LogCache struct {
	lines []string
}

func NewLogCache() *LogCache {
	return &LogCache{}
}

func (c *LogCache) Write(p []byte) (n int, err error) {
	s := string(p)
	for _, l := range strings.Split(s, "\n") {
		if l == "" {
			continue
		}
		if len(c.lines) > 100 {
			c.lines = c.lines[1:]
		}
		c.lines = append(c.lines, l)
	}
	return len(p), nil
}

func (c *LogCache) Lines() []string {
	return c.lines
}
