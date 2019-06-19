module github.com/travisjeffery/kube-leaderelection

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
	k8s.io/api v0.0.0-20190615205754-1d1b8b084b30
	k8s.io/apimachinery v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/client-go v9.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20190603182131-db7b694dc208 // indirect
	k8s.io/kubernetes v1.11.2
)

replace k8s.io/api => k8s.io/api v0.0.0-20190528110122-9ad12a4af326 // Pinned to release-1.12.9

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221084156-01f179d85dbc // Pinned to release-1.12.9

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190528110200-4f3abb12cae2 // Pinned to release-1.12.9

replace k8s.io/kubernetes => k8s.io/kubernetes v1.12.9

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190528110544-fa58353d80f3 // Pinned to release-1.12.9
