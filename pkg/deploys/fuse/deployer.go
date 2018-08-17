package fuse

import (
	"net/http"

	brokerapi "github.com/aerogear/managed-services-broker/pkg/broker"
	"github.com/aerogear/managed-services-broker/pkg/clients/openshift"
	k8sClient "github.com/operator-framework/operator-sdk/pkg/k8sclient"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	"github.com/pkg/errors"
	glog "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type FuseDeployer struct {
	id string
}

func NewDeployer(id string) *FuseDeployer {
	return &FuseDeployer{id: id}
}

func (fd *FuseDeployer) DoesDeploy(serviceID string) bool {
	return serviceID == "fuse-service-id"
}

func (fd *FuseDeployer) GetCatalogEntries() []*brokerapi.Service {
	glog.Infof("Getting fuse catalog entries")
	return getCatalogServicesObj()
}

func (fd *FuseDeployer) Deploy(id string, k8sclient kubernetes.Interface, osClientFactory *openshift.ClientFactory) (*brokerapi.CreateServiceInstanceResponse, error) {
	glog.Infof("Deploying fuse from deployer, id: %s", id)

	// Namespace
	ns, err := k8sclient.CoreV1().Namespaces().Create(getNamespaceObj("fuse-" + id))
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create namespace for fuse service")
	}

	namespace := ns.ObjectMeta.Name

	// ServiceAccount
	_, err = k8sclient.CoreV1().ServiceAccounts(namespace).Create(getServiceAccountObj())
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create service account for fuse service")
	}

	//Role
	_, err = k8sclient.RbacV1beta1().Roles(namespace).Create(getRoleObj())
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create role for fuse service")
	}

	//RoleBindings
	for _, sysRoleBinding := range getSystemRoleBindings(namespace) {
		_, err = k8sclient.RbacV1beta1().RoleBindings(namespace).Create(&sysRoleBinding)
		if err != nil {
			return &brokerapi.CreateServiceInstanceResponse{
				Code: http.StatusInternalServerError,
			}, errors.Wrapf(err, "failed to create rolebinding for %s", &sysRoleBinding.ObjectMeta.Name)
		}
	}

	// Fuse specific rolebindings
	_, err = k8sclient.RbacV1beta1().RoleBindings(namespace).Create(getInstallRoleBindingObj())
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create install role binding for fuse service")
	}

	authClient, err := osClientFactory.AuthClient()
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create an openshift authorization client")
	}

	_, err = authClient.RoleBindings(namespace).Create(getViewRoleBindingObj())
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create view role binding for fuse service")
	}

	_, err = authClient.RoleBindings(namespace).Create(getEditRoleBindingObj())
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create edit role binding for fuse service")
	}

	// ImageStream
	imageClient, err := osClientFactory.ImageStreamClient()
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create an openshift image stream client")
	}

	_, err = imageClient.ImageStreams(namespace).Create(getImageStreamObj())
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create image stream for fuse service")
	}

	// DeploymentConfig
	dcClient, err := osClientFactory.AppsClient()
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create an openshift deployment config client")
	}

	_, err = dcClient.DeploymentConfigs(namespace).Create(getDeploymentConfigObj())
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create deployment config for fuse service")
	}

	// CRD ?

	// CR
	fuseClient, _, err := k8sClient.GetResourceClient("syndesis.io/v1alpha1", "Syndesis", namespace)
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create fuse client")
	}
	glog.Infof("Created Fuse client: %v", fuseClient)

	fuseObj := k8sutil.UnstructuredFromRuntimeObject(getFuseObj())

	fuseCR, err := fuseClient.Create(fuseObj)
	if err != nil {
		return &brokerapi.CreateServiceInstanceResponse{
			Code: http.StatusInternalServerError,
		}, errors.Wrap(err, "failed to create a fuse custom resource")
	}

	glog.Infof("FUSE CR created: +%v", fuseCR)
	return &brokerapi.CreateServiceInstanceResponse{
		Code:         http.StatusAccepted,
		DashboardURL: "",
	}, nil
}

func (fd *FuseDeployer) LastOperation(instanceID string, k8sclient kubernetes.Interface, osclient *openshift.ClientFactory) (*brokerapi.LastOperationResponse, error) {
	return &brokerapi.LastOperationResponse{
		State:       brokerapi.StateSucceeded,
		Description: "deploying fuse",
	}, nil
}

func (fd *FuseDeployer) GetID() string {
	glog.Infof("getting fuse id from deployer: " + fd.id)
	return fd.id
}
