package build

import (
	"os"
	"strings"
)

// Shelly splits a string at spaces, taking into account shell quotes and {{template directives}}
func Shelly(s string, env ...map[string]any) []string {
	context := map[string]any{}
	for _, line := range os.Environ() {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}
		context[parts[0]] = parts[1]
	}
	for k, v := range Globals {
		context[k] = v
	}
	for _, arg := range env {
		for k, v := range arg {
			context[k] = v
		}
	}

	var out []string
	start := 0
	var quote rune = 0
	for i, ch := range s {
		switch ch {
		case '{':
			quote = ch
		case '}':
			quote = 0
		case '\'', '"', '`':
			if quote == ch {
				out = append(out, s[start:i])
				start = i + 1
				quote = 0
				continue
			}
			if quote > 0 {
				continue
			}
			quote = ch
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		case ' ', '\t', '\r', '\n':
			if quote > 0 {
				break
			}
			v := strings.TrimSpace(s[start:i])
			if len(v) > 0 {
				out = append(out, v)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
