# chart-doc-gen
Helm Chart Documentation Generator

```console
$ go run *.go \
    -d=./testdata/doc.yaml \
    -v=./testdata/values.yaml > ./testdata/README.md
```

## Installation

Download the pre-built binaries from release page and copy to your $PATH directory. If you are using Go modules, you can install like below:

```console
go get -u kmodules.xyz/chart-doc-gen@v0.2.8
```

## How does it work

`chart-doc-gen` takes a [doc.yaml](./testdata/doc.yaml) file and fills it with a values table auto generated from a [chart values file](./testdata/values.yaml). Then it renders to stdout a README.md file based on [readme template](./templates/readme.tpl).

`chart-doc-gen` walks a chart values file and generates a row for each leaf node in the values YAML document.
The description of each leaf node must be written above it in comments.
You can find an example generated [README.md](./testdata/README.md) from a [values file](./testdata/values.yaml).

Sometimes you may provide an object as default value for a parameter. To break out from the tree walk in that case,
add a line comment `+doc-gen:break` to the right of the parameter.

![+doc-gen:break example](./images/doc_gen_break.png "+doc-gen:break example")
![+doc-gen:break preview](./images/doc_gen_break_preview.png "+doc-gen:break preview")

You can also add an example for `--set key=value` command in the comments. To do so, add a line `# Example:`
and write the example commands in the comments below. The example lines will be broken by `<br >` in the
generated values table.

![values example](./images/values-example.png "Example in Description")
![values example preview](./images/values-example-preview.png "Preview Example in Description")

## Use with CI

You can use this tool in CI pipelines to ensure that your chart readme is up-to-date. You can use a Makefile with targets like below:

```console
.PHONY: gen
gen: gen-chart-doc

.PHONY: gen-chart-doc
gen-chart-doc:
	@echo "Generate chart docs"
	@chart-doc-gen -d=./testdata/doc.yaml -v=./testdata/values.yaml > ./testdata/README.md

.PHONY: verify
verify: verify-gen

.PHONY: verify-gen
verify-gen: gen fmt
	@if !(git diff --exit-code HEAD); then \
		echo "generated files are out of date, run make gen"; exit 1; \
	fi
```
