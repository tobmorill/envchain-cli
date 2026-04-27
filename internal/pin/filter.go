package pin

import "github.com/envchain/envchain-cli/internal/env"

// Filter returns only the entries from all whose key appears in the
// pinned key list for project. If the project has no pins (ErrNotFound)
// the original slice is returned unchanged so callers always get a
// usable result.
func (m *Manager) Filter(project string, all []env.Entry) ([]env.Entry, error) {
	keys, err := m.Get(project)
	if err == ErrNotFound {
		return all, nil
	}
	if err != nil {
		return nil, err
	}
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	out := make([]env.Entry, 0, len(keys))
	for _, e := range all {
		if _, ok := set[e.Key]; ok {
			out = append(out, e)
		}
	}
	return out, nil
}
