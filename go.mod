module infinibox-csi-driver

go 1.19

require (
	bou.ke/monkey v1.0.2
	github.com/container-storage-interface/spec v1.5.0
	github.com/containerd/containerd v1.5.16
	github.com/go-resty/resty/v2 v2.6.0
	github.com/golang/protobuf v1.5.2
	github.com/rexray/gocsi v1.2.2
	github.com/stretchr/testify v1.8.1
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
	k8s.io/api v0.24.9
	k8s.io/apimachinery v0.24.9
	k8s.io/client-go v0.24.9
	k8s.io/klog v1.0.0
	k8s.io/kubernetes v1.24.9
	k8s.io/mount-utils v0.24.9
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
)

require (
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/Microsoft/hcsshim v0.9.5 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/akutz/gosync v0.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/cgroups v1.0.4 // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful v2.9.5+incompatible // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/opencontainers/selinux v1.10.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	golang.org/x/tools v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221207170731-23e4bf6bdc37 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiserver v0.24.9 // indirect
	k8s.io/cloud-provider v0.0.0 // indirect
	k8s.io/component-base v0.24.9 // indirect
	k8s.io/component-helpers v0.24.9 // indirect
	k8s.io/klog/v2 v2.60.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220328201542-3ee0da9b0b42 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace k8s.io/api => k8s.io/api v0.24.9

replace k8s.io/apimachinery => k8s.io/apimachinery v0.24.10-rc.0

replace k8s.io/client-go => k8s.io/client-go v0.24.9

replace k8s.io/mount-utils => k8s.io/mount-utils v0.24.10-rc.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.24.9

replace k8s.io/apiserver => k8s.io/apiserver v0.24.9

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.24.9

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.24.9

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.24.9

replace k8s.io/code-generator => k8s.io/code-generator v0.24.10-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.24.9

replace k8s.io/component-helpers => k8s.io/component-helpers v0.24.9

replace k8s.io/controller-manager => k8s.io/controller-manager v0.24.9

replace k8s.io/cri-api => k8s.io/cri-api v0.24.10-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.24.9

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.24.9

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.24.9

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.24.9

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.24.9

replace k8s.io/kubectl => k8s.io/kubectl v0.24.9

replace k8s.io/kubelet => k8s.io/kubelet v0.24.9

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.24.9

replace k8s.io/metrics => k8s.io/metrics v0.24.9

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.24.9

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.24.9

replace k8s.io/sample-controller => k8s.io/sample-controller v0.24.9

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.24.9
