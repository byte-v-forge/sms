package herosms

func heroSMSServiceCandidates(applicationKey string) []string {
	serviceKey := normalizeHeroSMSServiceKey(applicationKey)
	if serviceKey == "" {
		return []string{""}
	}
	return []string{serviceKey}
}
