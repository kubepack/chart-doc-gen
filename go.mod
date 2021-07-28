module kubepack.dev/chart-doc-gen

go 1.14

require (
	github.com/go-errors/errors v1.0.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.21.1
	kmodules.xyz/client-go dd0503cf99cf3b6abb635d8945a8d7d8fed901d9
	kmodules.xyz/custom-resources 83db827677cf5651491478fa85707d62416cf679
	kmodules.xyz/resource-metadata dcc1abc08aa00646b9474f7702b45c798b3ce66c
	kmodules.xyz/webhook-runtime e489faf01981d2f3afa671989388c7b6f22b6baa
	sigs.k8s.io/kustomize/kyaml v0.1.12
)

replace github.com/satori/go.uuid => github.com/gofrs/uuid v4.0.0+incompatible

replace helm.sh/helm/v3 => github.com/kubepack/helm/v3 v3.6.1-0.20210518225915-c3e0ce48dd1b

replace k8s.io/apiserver => github.com/kmodules/apiserver v0.21.2-0.20210716212718-83e5493ac170
