module github.com/spring-cloud-incubator/mononoke

go 1.13

require (
	github.com/google/go-cmp v0.4.0
	github.com/google/go-containerregistry v0.0.0-20200304201134-fcc8ea80e26f
	github.com/projectriff/system v0.5.0
	golang.org/x/tools v0.0.0-20200306191617-51e69f71924f // indirect
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
	sigs.k8s.io/controller-runtime v0.5.0
	sigs.k8s.io/controller-tools v0.2.4
)

replace github.com/Azure/go-autorest v10.15.5+incompatible => github.com/Azure/go-autorest/autorest v0.9.0
