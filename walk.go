/*
Copyright The Kubepack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strings"

	"kubepack.dev/chart-doc-gen/walk"

	"sigs.k8s.io/kustomize/kyaml/openapi"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// GenerateValuesTable produces a table of parameters from a chart values file
func GenerateValuesTable(from *yaml.RNode) ([][]string, error) {
	p := printer{
		rows: [][]string{},
	}
	_, err := walk.Walker{
		Source:             from,
		Visitor:            &p,
		VisitKeysAsScalars: true,
	}.Walk()
	return p.rows, err
}

type printer struct {
	rows [][]string
}

func (c *printer) VisitMap(s *yaml.RNode, _ *openapi.ResourceSchema) (*yaml.RNode, error) {
	return s, nil
}

func (c *printer) VisitScalar(s *yaml.RNode, _ *openapi.ResourceSchema) (*yaml.RNode, error) {
	return s, nil
}

func (c *printer) VisitList(s *yaml.RNode, _ *openapi.ResourceSchema, _ walk.ListKind) (*yaml.RNode, error) {
	return s, nil
}

func (c *printer) VisitLeaf(key *yaml.RNode, value *yaml.RNode, path string, _ *openapi.ResourceSchema) (*yaml.RNode, error) {
	desc, example := ParseComment(key.YNode().HeadComment)
	c.rows = append(c.rows, []string{
		path,
		desc,
		PrintValue(value),
		example,
	})
	return key, nil
}

func ParseComment(s string) (string, string) {
	lines := strings.Split(s, "\n")
	var desc []string
	var example []string
	idx := 0
	for ; idx < len(lines); idx++ {
		line := walk.CommentValue(lines[idx])
		if line == "Example:" || line == "Example" {
			break
		}
		if line != "" {
			desc = append(desc, line)
		}
	}

	for idx++; idx < len(lines); idx++ {
		line := walk.CommentValue(lines[idx])
		if line != "" {
			example = append(example, fmt.Sprintf("`%s`", line))
		}
	}
	eg := strings.Join(example, " <br> ")
	d := strings.Join(desc, " ")
	if len(eg) > 0 {
		d = d + " <br> Example: <br> " + eg
	}
	return d, eg
}

func PrintValue(node *yaml.RNode) string {
	if node.YNode().Kind == yaml.MappingNode || node.YNode().Kind == yaml.SequenceNode {
		if !yaml.IsMissingOrNull(node) {
			data, err := node.MarshalJSON()
			if err != nil {
				panic(err)
			}
			return string(data)
		}
	}
	return strings.TrimSpace(node.MustString())
}
