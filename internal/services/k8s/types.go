package k8s

import (
	"context"
	"fmt"
	"strconv"

	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func clusterAutoscalerConfigFlatten(cluster *k8s.Cluster) []map[string]any {
	autoscalerConfig := map[string]any{}
	autoscalerConfig["disable_scale_down"] = cluster.AutoscalerConfig.ScaleDownDisabled
	autoscalerConfig["scale_down_delay_after_add"] = cluster.AutoscalerConfig.ScaleDownDelayAfterAdd
	autoscalerConfig["scale_down_unneeded_time"] = cluster.AutoscalerConfig.ScaleDownUnneededTime
	autoscalerConfig["estimator"] = cluster.AutoscalerConfig.Estimator
	autoscalerConfig["expander"] = cluster.AutoscalerConfig.Expander
	autoscalerConfig["ignore_daemonsets_utilization"] = cluster.AutoscalerConfig.IgnoreDaemonsetsUtilization
	autoscalerConfig["balance_similar_node_groups"] = cluster.AutoscalerConfig.BalanceSimilarNodeGroups
	autoscalerConfig["expendable_pods_priority_cutoff"] = cluster.AutoscalerConfig.ExpendablePodsPriorityCutoff

	// need to convert a f32 to f64 without precision loss
	thresholdF64, err := strconv.ParseFloat(fmt.Sprintf("%f", cluster.AutoscalerConfig.ScaleDownUtilizationThreshold), 64)
	if err != nil {
		// should never happen
		return nil
	}

	autoscalerConfig["scale_down_utilization_threshold"] = thresholdF64
	autoscalerConfig["max_graceful_termination_sec"] = cluster.AutoscalerConfig.MaxGracefulTerminationSec

	return []map[string]any{autoscalerConfig}
}

func clusterOpenIDConnectConfigFlatten(cluster *k8s.Cluster) []map[string]any {
	openIDConnectConfig := map[string]any{}
	openIDConnectConfig["issuer_url"] = cluster.OpenIDConnectConfig.IssuerURL
	openIDConnectConfig["client_id"] = cluster.OpenIDConnectConfig.ClientID
	openIDConnectConfig["username_claim"] = cluster.OpenIDConnectConfig.UsernameClaim
	openIDConnectConfig["username_prefix"] = cluster.OpenIDConnectConfig.UsernamePrefix
	openIDConnectConfig["groups_claim"] = cluster.OpenIDConnectConfig.GroupsClaim
	openIDConnectConfig["groups_prefix"] = cluster.OpenIDConnectConfig.GroupsPrefix
	openIDConnectConfig["required_claim"] = cluster.OpenIDConnectConfig.RequiredClaim

	return []map[string]any{openIDConnectConfig}
}

func clusterAutoUpgradeFlatten(cluster *k8s.Cluster) []map[string]any {
	autoUpgrade := map[string]any{}
	autoUpgrade["enable"] = cluster.AutoUpgrade.Enabled
	autoUpgrade["maintenance_window_start_hour"] = cluster.AutoUpgrade.MaintenanceWindow.StartHour
	autoUpgrade["maintenance_window_day"] = cluster.AutoUpgrade.MaintenanceWindow.Day

	return []map[string]any{autoUpgrade}
}

func poolUpgradePolicyFlatten(pool *k8s.Pool) []map[string]any {
	upgradePolicy := map[string]any{}
	if pool.UpgradePolicy != nil {
		upgradePolicy["max_surge"] = pool.UpgradePolicy.MaxSurge
		upgradePolicy["max_unavailable"] = pool.UpgradePolicy.MaxUnavailable
	}

	return []map[string]any{upgradePolicy}
}

func expandKubeletArgs(args any) map[string]string {
	kubeletArgs := map[string]string{}

	for key, value := range args.(map[string]any) {
		kubeletArgs[key] = value.(string)
	}

	return kubeletArgs
}

func flattenKubeletArgs(args map[string]string) map[string]any {
	kubeletArgs := map[string]any{}

	for key, value := range args {
		kubeletArgs[key] = value
	}

	return kubeletArgs
}

func flattenKubeconfig(ctx context.Context, k8sAPI *k8s.API, region scw.Region, clusterID string) (map[string]any, error) {
	kubeconfig, err := k8sAPI.GetClusterKubeConfig(&k8s.GetClusterKubeConfigRequest{
		Region:    region,
		ClusterID: clusterID,
	}, scw.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	kubeconfigServer, err := kubeconfig.GetServer()
	if err != nil {
		return nil, err
	}

	kubeconfigCa, err := kubeconfig.GetCertificateAuthorityData()
	if err != nil {
		return nil, err
	}

	kubeconfigToken, err := kubeconfig.GetToken()
	if err != nil {
		return nil, err
	}

	kubeconf := map[string]any{}
	kubeconf["config_file"] = string(kubeconfig.GetRaw())
	kubeconf["host"] = kubeconfigServer
	kubeconf["cluster_ca_certificate"] = kubeconfigCa
	kubeconf["token"] = kubeconfigToken

	return kubeconf, nil
}
