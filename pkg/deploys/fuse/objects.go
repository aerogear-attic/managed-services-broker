package fuse

import (
	brokerapi "github.com/aerogear/managed-services-broker/pkg/broker"
	"github.com/aerogear/managed-services-broker/pkg/deploys/fuse/pkg/apis/syndesis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Fuse plan
func getCatalogServicesObj() []*brokerapi.Service {
	return []*brokerapi.Service{
		{
			Name:        "fuse",
			ID:          "fuse-service-id",
			Description: "fuse",
			Metadata:    map[string]string{"serviceName": "fuse", "serviceType": "fuse"},
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					Name:        "default-fuse",
					ID:          "default-fuse",
					Description: "default fuse plan",
					Free:        true,
					Schemas: &brokerapi.Schemas{
						ServiceBinding: &brokerapi.ServiceBindingSchema{
							Create: &brokerapi.RequestResponseSchema{},
						},
						ServiceInstance: &brokerapi.ServiceInstanceSchema{
							Create: &brokerapi.InputParametersSchema{},
						},
					},
				},
			},
		},
	}
}

// Fuse Custom Resource
func getFuseObj(serviceInstanceId string, userNamespace string) *v1alpha1.Syndesis {
	demoData := false
	deployIntegrations := true
	limit := 1
	stateCheckInterval := 60

	return &v1alpha1.Syndesis{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Syndesis",
			APIVersion: "syndesis.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "fuse-" + "-" + serviceInstanceId,
			Annotations: map[string]string{},
		},
		Spec: v1alpha1.SyndesisSpec{
			SarNamespace:         userNamespace,
			DemoData:             &demoData,
			DeployIntegrations:   &deployIntegrations,
			ImageStreamNamespace: "",
			Integration: v1alpha1.IntegrationSpec{
				Limit:              &limit,
				StateCheckInterval: &stateCheckInterval,
			},
			Registry: "docker.io",
			Components: v1alpha1.ComponentsSpec{
				Db: v1alpha1.DbConfiguration{
					Resources: v1alpha1.ResourcesWithVolume{
						ResourceRequirements: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"memory": *resource.NewQuantity(255*1024*1024, resource.BinarySI),
							},
						},
						VolumeCapacity: "1Gi",
					},
					User:                 "syndesis",
					Database:             "syndesis",
					ImageStreamNamespace: "openshift",
				},
				Prometheus: v1alpha1.PrometheusConfiguration{
					Resources: v1alpha1.ResourcesWithVolume{
						ResourceRequirements: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"memory": *resource.NewQuantity(512*1024*1024, resource.BinarySI),
							},
						},
						VolumeCapacity: "1Gi",
					},
				},
				Server: v1alpha1.ServerConfiguration{
					Resources: v1alpha1.Resources{
						ResourceRequirements: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"memory": *resource.NewQuantity(800*1024*1024, resource.BinarySI),
							},
						},
					},
				},
				Meta: v1alpha1.MetaConfiguration{
					Resources: v1alpha1.ResourcesWithVolume{
						ResourceRequirements: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"memory": *resource.NewQuantity(512*1024*1024, resource.BinarySI),
							},
						},
						VolumeCapacity: "1Gi",
					},
				},
			},
		},
	}
}
