package mongodb

import (
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
)

func flattenPublicNetwork(endpoints []*mongodb.Endpoint) (any, bool) {
	publicFlat := []map[string]any(nil)

	for _, endpoint := range endpoints {
		if endpoint.Public == nil {
			continue
		}

		publicFlat = append(publicFlat, map[string]any{
			"id":         endpoint.ID,
			"port":       endpoint.Port,
			"dns_record": endpoint.DNSRecords[0],
		})

		break
	}

	return publicFlat, len(publicFlat) != 0
}

func flattenPrivateNetwork(endpoints []*mongodb.Endpoint) (any, bool) {
	privateFlat := []map[string]any(nil)

	for _, endpoint := range endpoints {
		if endpoint.PrivateNetwork == nil {
			continue
		}

		privateFlat = append(privateFlat, map[string]any{
			"pn_id":       endpoint.PrivateNetwork.PrivateNetworkID,
			"id":          endpoint.ID,
			"port":        endpoint.Port,
			"dns_records": endpoint.DNSRecords,
			"ips":         endpoint.IPs,
		})

		break
	}

	return privateFlat, len(privateFlat) != 0
}
