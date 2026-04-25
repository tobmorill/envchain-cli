// Package tags provides project tagging functionality for envchain-cli.
//
// Tags allow users to categorise and group projects by assigning
// one or more short labels (e.g. "work", "staging", "personal").
// Tags are stored alongside the project metadata and can be used
// to filter results when searching or listing chains.
//
// Usage:
//
//	manager := tags.New(store)
//
//	// Assign tags to a project
//	if err := manager.Set("myproject", []string{"work", "staging"}); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Retrieve tags for a project
//	t, err := manager.Get("myproject")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(t) // [staging work]
//
//	// Find all projects that carry a given tag
//	projects, err := manager.FindByTag("work")
//	if err != nil {
//	    log.Fatal(err)
//	}
package tags
