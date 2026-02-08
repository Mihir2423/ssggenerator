package markdown

import "strings"

func ToHTML(input []byte) string {
	lines := strings.Split(string(input), "\n")
	var out strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "## "):
			out.WriteString("<h2>")
			out.WriteString(strings.TrimPrefix(line, "## "))
			out.WriteString("</h2>")
		case strings.HasPrefix(line, "# "):
			out.WriteString("<h1>")
			out.WriteString(strings.TrimPrefix(line, "# "))
			out.WriteString("</h1>")
		default:
			out.WriteString("<p>")
			out.WriteString(line)
			out.WriteString("</p>")
		}
	}
	return out.String()
}
