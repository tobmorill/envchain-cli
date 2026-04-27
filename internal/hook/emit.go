package hook

import (
	"fmt"
	"strings"
)

// EmitScript returns a shell snippet that executes the hook command for the
// given project and phase. If no hook is registered the returned string is
// empty. The shell parameter should be "bash", "zsh", or "fish".
func EmitScript(m *Manager, project string, phase Phase, shell string) (string, error) {
	h, ok, err := m.Get(project, phase)
	if err != nil {
		return "", fmt.Errorf("hook emit: %w", err)
	}
	if !ok {
		return "", nil
	}
	return buildSnippet(h.Command, phase, shell), nil
}

func buildSnippet(command string, phase Phase, shell string) string {
	var sb strings.Builder
	switch shell {
	case "fish":
		fmt.Fprintf(&sb, "# envchain hook: %s\n", phase)
		fmt.Fprintf(&sb, "%s\n", command)
	default: // bash / zsh / posix
		fmt.Fprintf(&sb, "# envchain hook: %s\n", phase)
		fmt.Fprintf(&sb, "( %s )\n", command)
	}
	return sb.String()
}
