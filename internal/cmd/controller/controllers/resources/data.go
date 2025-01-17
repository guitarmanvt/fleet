package resources

import (
	fleetgroup "github.com/rancher/fleet/pkg/apis/fleet.cattle.io"
	fleet "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"

	"github.com/rancher/wrangler/v2/pkg/apply"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	BundleDeploymentClusterRole = "fleet-bundle-deployment"
	ContentClusterRole          = "fleet-content"
)

// ApplyBootstrapResources creates the cluster roles, system namespace and system registration namespace
func ApplyBootstrapResources(systemNamespace, systemRegistrationNamespace string, apply apply.Apply) error {
	return apply.ApplyObjects(
		// used by request-* service accounts from agents
		&rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: BundleDeploymentClusterRole,
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get", "list", "watch"},
					APIGroups: []string{fleetgroup.GroupName},
					Resources: []string{fleet.BundleDeploymentResourceName},
				},
				{
					Verbs:     []string{"update"},
					APIGroups: []string{fleetgroup.GroupName},
					Resources: []string{fleet.BundleDeploymentResourceName + "/status"},
				},
			},
		},
		// used by request-* service accounts from agents
		&rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: ContentClusterRole,
			},
			Rules: []rbacv1.PolicyRule{
				{
					Verbs:     []string{"get"},
					APIGroups: []string{fleetgroup.GroupName},
					Resources: []string{fleet.ContentResourceName},
				},
			},
		},
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: systemNamespace,
			},
		},
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: systemRegistrationNamespace,
			},
		},
	)
}
