package fuse

import (
	brokerapi "github.com/aerogear/managed-services-broker/pkg/broker"
	corev1 "k8s.io/api/core/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getCatalogServicesObj() []*brokerapi.Service {
	return []*brokerapi.Service{
		{
			Name:        "fuse",
			ID:          "fuse-service-id",
			Description: "fuse",
			Metadata:    map[string]string{"serviceName": "fuse", "serviceType": "fuse"},
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					Name:        "default",
					ID:          "default",
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

func getNamespaceObj(id string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
	}
}

func getServiceAccountObj() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "syndesis-operator",
			Labels: map[string]string{
				"app":                   "syndesis",
				"syndesis.io/app":       "syndesis",
				"syndesis.io/type":      "operator",
				"syndesis.io/component": "syndesis-operator",
			},
		},
	}
}

func getRoleObj() *rbacv1beta1.Role {
	return &rbacv1beta1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: "syndesis-operator",
			Labels: map[string]string{
				"app":                   "syndesis",
				"syndesis.io/app":       "syndesis",
				"syndesis.io/type":      "operator",
				"syndesis.io/component": "syndesis-operator",
			},
		},
		Rules: []rbacv1beta1.PolicyRule{
			{
				APIGroups: []string{"syndesis.io"},
				Resources: []string{"syndesises", "syndesises/finalizers"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "services", "endpoints", "persistentvolumeclaims", "configmaps", "secrets", "serviceaccounts"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"get", "list"},
			},
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"rolebindings"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
			{
				APIGroups: []string{"template.openshift.io"},
				Resources: []string{"processedtemplates"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
			{
				APIGroups: []string{"image.openshift.io"},
				Resources: []string{"imagestreams"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
			{
				APIGroups: []string{"apps.openshift.io"},
				Resources: []string{"deploymentconfigs"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
			{
				APIGroups: []string{"build.openshift.io"},
				Resources: []string{"buildconfigs"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
			{
				APIGroups: []string{"authorization.openshift.io"},
				Resources: []string{"rolebindings"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
			{
				APIGroups: []string{"route.openshift.io"},
				Resources: []string{"routes", "routes/custom-host"},
				Verbs:     []string{"get", "list", "create", "update", "delete", "deletecollection", "watch"},
			},
		},
	}
}

func getInstallRoleBindingObj() *rbacv1beta1.RoleBinding {
	return &rbacv1beta1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "syndesis-operator:install",
			Labels: map[string]string{
				"app":                   "syndesis",
				"syndesis.io/app":       "syndesis",
				"syndesis.io/type":      "operator",
				"syndesis.io/component": "syndesis-operator",
			},
		},
		Subjects: []rbacv1beta1.Subject{
			{
				Kind: "ServiceAccount",
				Name: "syndesis-operator",
			},
		},
		RoleRef: rbacv1beta1.RoleRef{
			Kind:     "Role",
			Name:     "syndesis-operator",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}
