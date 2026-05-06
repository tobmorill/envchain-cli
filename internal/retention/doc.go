// Package retention provides per-project data retention policies for
// envchain-cli. A Policy specifies how long historical snapshots and
// version records should be kept (MaxAge) and the maximum number of
// versions to retain (MaxVersions).
//
// Policies are stored in the shared envchain store under a namespaced
// key and can be retrieved, updated, or deleted at any time.
//
// Example usage:
//
//	pol := retention.Policy{
//		Project:     "myapp",
//		MaxAge:      7 * 24 * time.Hour,
//		MaxVersions: 20,
//	}
//	if err := mgr.Set(pol); err != nil { ... }
package retention
