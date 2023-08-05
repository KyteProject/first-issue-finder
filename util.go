package main

import (
	"strings"

	"github.com/google/go-github/v53/github"
)

// safeDereferenceLabels dereferences a slice of labels and returns a string of comma-separated labels.
func safeDereferenceLabels(ls []*github.Label) string {
	labels := make([]string, 0, len(ls))
	for _, label := range ls {
		if label != nil {
			labels = append(labels, label.GetName())
		}
	}
	return strings.Join(labels, ", ")
}
