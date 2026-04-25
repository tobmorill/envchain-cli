// Package template manages named collections of environment variable keys
// that can be used as reusable blueprints across envchain projects.
//
// A Template stores only the key names (not values), allowing users to define
// a standard set of required variables for a class of project (e.g. "django",
// "node-api") and later apply them when creating or editing a chain.
//
// Templates are persisted in the same key-value store used by other envchain
// components, namespaced under the "template::" prefix to avoid collisions.
package template
