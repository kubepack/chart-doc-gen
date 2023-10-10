// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package walk

import (
	"sigs.k8s.io/kustomize/kyaml/fieldmeta"
	"sigs.k8s.io/kustomize/kyaml/openapi"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"sigs.k8s.io/kustomize/kyaml/yaml/schema"
)

// Walker walks the Source RNode and modifies the RNode provided to GrepFilter.
type Walker struct {
	// Visitor is invoked by GrepFilter
	Visitor

	Schema *openapi.ResourceSchema

	// Source is the RNode to walk.  All Source fields and associative list elements
	// will be visited.
	Source *yaml.RNode

	// Path is the field path to the current Source Node.
	Path []string

	// InferAssociativeLists if set to true will infer merge strategies for
	// fields which it doesn't have the schema based on the fields in the
	// list elements.
	InferAssociativeLists bool

	// VisitKeysAsScalars if true will call VisitScalar on map entry keys,
	// providing nil as the OpenAPI schema.
	VisitKeysAsScalars bool
}

func (l Walker) Kind() yaml.Kind {
	if !yaml.IsMissingOrNull(l.Source) {
		return l.Source.YNode().Kind
	}
	return 0
}

// GrepFilter implements yaml.GrepFilter
func (l Walker) Walk() (*yaml.RNode, error) {
	l.Schema = l.GetSchema()

	// invoke the handler for the corresponding node type
	switch l.Kind() {
	case yaml.MappingNode:
		if err := yaml.ErrorIfAnyInvalidAndNonNull(yaml.MappingNode, l.Source); err != nil {
			return nil, err
		}
		return l.walkMap()
	case yaml.SequenceNode:
		if err := yaml.ErrorIfAnyInvalidAndNonNull(yaml.SequenceNode, l.Source); err != nil {
			return nil, err
		}
		if schema.IsAssociative(l.Schema, []*yaml.RNode{l.Source}, l.InferAssociativeLists) {
			return l.walkAssociativeSequence()
		}
		return l.walkNonAssociativeSequence()

	case yaml.ScalarNode:
		if err := yaml.ErrorIfAnyInvalidAndNonNull(yaml.ScalarNode, l.Source); err != nil {
			return nil, err
		}
		return l.walkScalar()
	case 0:
		// walk empty nodes as maps
		return l.walkMap()
	default:
		return nil, nil
	}
}

func (l Walker) GetSchema() *openapi.ResourceSchema {
	r := l.Source
	if yaml.IsMissingOrNull(r) {
		return nil
	}

	fm := fieldmeta.FieldMeta{}
	if err := fm.Read(r); err == nil && !fm.IsEmpty() {
		// per-field schema, this is fine
		if fm.Schema.Ref.String() != "" {
			// resolve the reference
			s, err := openapi.Resolve(&fm.Schema.Ref, &fm.Schema)
			if err == nil && s != nil {
				fm.Schema = *s
			}
		}
		return &openapi.ResourceSchema{Schema: &fm.Schema}
	}

	if l.Schema != nil {
		return l.Schema
	}

	r = l.Source
	if yaml.IsMissingOrNull(r) {
		return nil
	}

	m, _ := r.GetMeta()
	if m.Kind == "" || m.APIVersion == "" {
		return nil
	}

	s := openapi.SchemaForResourceType(yaml.TypeMeta{Kind: m.Kind, APIVersion: m.APIVersion})
	if s != nil {
		return s
	}
	return nil
}
