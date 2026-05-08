// Package embargo implements time-window based access restrictions for
// envchain project chains.
//
// An embargo window specifies a daily UTC time range during which a project's
// environment variables may be loaded. Attempts to access the chain outside
// that window are rejected with ErrEmbargoActive.
//
// Windows that span midnight (e.g. 22:00–06:00) are supported: when the start
// hour is greater than the end hour the window is treated as crossing the
// day boundary.
//
// Usage:
//
//	st, _ := store.New(path)
//	m := embargo.New(st)
//	m.Set("myproject", embargo.Window{StartHour: 9, EndHour: 17})
//	if err := m.Check("myproject"); err != nil {
//	    // access denied
//	}
package embargo
