package openshift

import (
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	authv1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	imagev1 "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	"k8s.io/client-go/rest"
)

func NewClientFactory(cfg *rest.Config) *ClientFactory {
	return &ClientFactory{cfg: cfg}
}

type ClientFactory struct {
	cfg *rest.Config
}

func (c *ClientFactory) BuildClient() (*buildv1.BuildV1Client, error) {
	return buildv1.NewForConfig(c.cfg)
}

func (c *ClientFactory) AuthClient() (*authv1.AuthorizationV1Client, error) {
	return authv1.NewForConfig(c.cfg)
}

func (c *ClientFactory) ImageStreamClient() (*imagev1.ImageV1Client, error) {
	return imagev1.NewForConfig(c.cfg)
}

func (c *ClientFactory) AppsClient() (*appsv1.AppsV1Client, error) {
	return appsv1.NewForConfig(c.cfg)
}
