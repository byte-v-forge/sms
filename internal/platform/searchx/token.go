package searchx

import "strings"

type Candidate struct {
	Key  string
	Name string
}

func Token(value string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
		}
	}
	return out.String()
}

func ContainsToken(candidate string, query string) bool {
	candidateToken := Token(candidate)
	queryToken := Token(query)
	return candidateToken != "" && queryToken != "" && (candidateToken == queryToken || strings.Contains(candidateToken, queryToken))
}

func MatchKey(query string, candidates []Candidate) string {
	queryToken := Token(query)
	if queryToken == "" {
		return ""
	}
	if key := exactKeyMatch(queryToken, candidates); key != "" {
		return key
	}
	if key := exactNameMatch(queryToken, candidates); key != "" {
		return key
	}
	return uniquePartialMatch(queryToken, candidates)
}

func exactKeyMatch(queryToken string, candidates []Candidate) string {
	for _, candidate := range candidates {
		if Token(candidate.Key) == queryToken {
			return strings.TrimSpace(candidate.Key)
		}
	}
	return ""
}

func exactNameMatch(queryToken string, candidates []Candidate) string {
	for _, candidate := range candidates {
		if Token(candidate.Name) == queryToken {
			return strings.TrimSpace(candidate.Key)
		}
	}
	return ""
}

func uniquePartialMatch(queryToken string, candidates []Candidate) string {
	matched := ""
	for _, candidate := range candidates {
		if !candidateMatches(queryToken, candidate) {
			continue
		}
		key := strings.TrimSpace(candidate.Key)
		if key == "" {
			continue
		}
		if matched != "" && matched != key {
			return ""
		}
		matched = key
	}
	return matched
}

func candidateMatches(queryToken string, candidate Candidate) bool {
	for _, value := range []string{candidate.Key, candidate.Name} {
		token := Token(value)
		if token != "" && strings.Contains(token, queryToken) {
			return true
		}
	}
	return false
}
