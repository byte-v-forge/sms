package herosms

func purchaseInfoForService(response purchaseInfoResponse, serviceKey string) purchaseInfo {
	serviceKey = normalizeHeroSMSServiceKey(serviceKey)
	if info, ok := response.Data[serviceKey]; ok {
		return info
	}
	for key, info := range response.Data {
		if normalizeHeroSMSServiceKey(key) == serviceKey {
			return info
		}
	}
	return purchaseInfo{}
}
