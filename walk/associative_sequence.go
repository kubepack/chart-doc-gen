// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package walk

import (
	"strings"

	"github.com/go-errors/errors"
	"sigs.k8s.io/kustomize/kyaml/openapi"
	"sigs.k8s.io/kustomize/kyaml/sets"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (l *Walker) walkAssociativeSequence() (*yaml.RNode, error) {
	// may require initializing the dest node
	dest, err := l.VisitList(l.Source, l.Schema, AssociativeList)
	if dest == nil || err != nil {
		return nil, err
	}
	l.Source = dest

	var key string
	if l.Schema != nil {
		_, key = l.Schema.PatchStrategyAndKey()
	}
	if key == "" { // no key from the schema, try to infer one
		// find the list of elements we need to recursively walk
		key, err = l.elementKey()
		if err != nil {
			return nil, err
		}
	}

	values := l.elementValues(key)

	// recursively set the elements in the list
	var s *openapi.ResourceSchema
	if l.Schema != nil {
		s = l.Schema.Elements()
	}
	for _, value := range values {
		val, err := Walker{
			VisitKeysAsScalars:    l.VisitKeysAsScalars,
			InferAssociativeLists: l.InferAssociativeLists,
			Visitor:               l,
			Schema:                s,
			Source:                l.elementValue(key, value),
		}.Walk()
		if err != nil {
			return nil, err
		}

		if yaml.IsMissingOrNull(val) {
			_, err = dest.Pipe(yaml.ElementSetter{Keys: []string{key}, Values: []string{value}})
			if err != nil {
				return nil, err
			}
			continue
		}

		if val.Field(key) == nil {
			// make sure the key is set on the field
			_, err = val.Pipe(yaml.SetField(key, yaml.NewScalarRNode(value)))
			if err != nil {
				return nil, err
			}
		}

		// this handles empty and non-empty values
		_, err = dest.Pipe(yaml.ElementSetter{Element: val.YNode(), Keys: []string{key}, Values: []string{value}})
		if err != nil {
			return nil, err
		}
	}
	// field is empty
	if yaml.IsMissingOrNull(dest) {
		return nil, nil
	}
	return dest, nil
}

// elementKey returns the merge key to use for the associative list
func (l Walker) elementKey() (string, error) {
	var key string
	if l.Source != nil && len(l.Source.Content()) > 0 {
		newKey := l.Source.GetAssociativeKey()
		if key != "" && key != newKey {
			return "", errors.Errorf(
				"conflicting merge keys [%s,%s] for field %s",
				key, newKey, strings.Join(l.Path, "."))
		}
		key = newKey
	}
	if key == "" {
		return "", errors.Errorf("no merge key found for field %s",
			strings.Join(l.Path, "."))
	}
	return key, nil
}

// elementValues returns a slice containing all values for the field across all elements
// from all sources.
// Return value slice is ordered using the original ordering from the elements, where
// elements missing from earlier sources appear later.
func (l Walker) elementValues(key string) []string {
	// use slice to to keep elements in the original order
	// dest node must be first
	var returnValues []string
	seen := sets.String{}
	if l.Source == nil {
		return nil
	}

	// add the value of the field for each element
	// don't check error, we know this is a list node
	values, _ := l.Source.ElementValues(key)
	for _, s := range values {
		if seen.Has(s) {
			continue
		}
		returnValues = append(returnValues, s)
		seen.Insert(s)
	}
	return returnValues
}

// fieldValue returns a slice containing each source's value for fieldName
func (l Walker) elementValue(key, value string) *yaml.RNode {
	if l.Source == nil {
		return nil
	}
	return l.Source.Element(key, value)
}
