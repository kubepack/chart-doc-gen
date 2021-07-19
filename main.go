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
	"bytes"
	"strings"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"kubepack.dev/chart-doc-gen/api"
	"kubepack.dev/chart-doc-gen/templates"

	"github.com/olekukonko/tablewriter"
	flag "github.com/spf13/pflag"
	y2 "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

var (
	docFile    *string = flag.StringP("doc", "d", "", "Path to a project's doc.{json|yaml} info file")
	valuesFile *string = flag.StringP("values", "v", "", "Path to chart values file")
	tplFile    *string = flag.StringP("template", "t", "readme2.tpl", "Path to a doc template file")
)

func main() {
	flag.Parse()

	f, err := os.Open(*docFile)
	if err != nil {
		panic(err)
	}
	reader := y2.NewYAMLOrJSONDecoder(f, 2048)
	var doc api.DocInfo
	err = reader.Decode(&doc)
	if err != nil && err != io.EOF {
		panic(err)
	}

	data, err := ioutil.ReadFile(*valuesFile)
	if err != nil {
		panic(err)
	}
	obj, err := yaml.Parse(string(data))
	if err == nil {
		rows, err := GenerateValuesTable(obj)
		if err != nil {
			panic(err)
		}

		var params [][]string
		for _, row := range rows {
			params = append(params, []string{
				row[0],
				row[1],
				fmt.Sprintf(
					"<code>%s</code>",  // use a html code block instead of backtics so the whole block get highlighted
					strings.ReplaceAll(  // replace all newlines, they generate new table columns with tablewriter
						strings.ReplaceAll(row[2], "|", "&#124;"),  // replace all pipe symbols with their ACSII representation, because they break the markdown table
						"\n",
						"&#13;&#10;",
					),
				),
			})
		}

		var buf bytes.Buffer
		table := tablewriter.NewWriter(&buf)
		table.SetHeader([]string{"Parameter", "Description", "Default"})
		table.SetAutoFormatHeaders(false)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetAutoWrapText(false)
		table.SetCenterSeparator("|")
		table.AppendBulk(params) // Add Bulk Data
		table.Render()

		doc.Chart.Values = buf.String()

		if doc.Chart.ValuesExample == "" || strings.HasPrefix(doc.Chart.ValuesExample, "-- generate from values file --") {
			for _, row := range rows {
				if row[2] != "" &&
					row[2] != `""` &&
					row[2] != "{}" &&
					row[2] != "[]" &&
					row[2] != "true" &&
					row[2] != "false" &&
					row[2] != "not-ca-cert" {
					doc.Chart.ValuesExample = fmt.Sprintf("%v=%v", row[0], row[2])
					break
				}
			}
		}
	} else if err == io.EOF {
		doc.Chart.Values = ""
		doc.Chart.ValuesExample = ""
	} else {
		panic(err)
	}

	tplReadme, err := ioutil.ReadFile(*tplFile)
	if err != nil {
		if os.IsNotExist(err) {
			tplReadme, err = templates.Asset("readme.tpl")
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	tmpl, err := template.New("readme").Parse(string(tplReadme))
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, doc)
	if err != nil {
		panic(err)
	}
}
