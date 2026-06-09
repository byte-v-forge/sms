package searchx

import "strings"

type Candidate struct {
	Key  string
	Name string
}

func Token(value string) string {
	var out strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if isSearchRune(r) {
			out.WriteRune(r)
		}
	}
	return out.String()
}

func Terms(value string) []string {
	var terms []string
	var current strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if isSearchRune(r) {
			current.WriteRune(r)
			continue
		}
		terms = appendSearchTerm(terms, &current)
	}
	return appendSearchTerm(terms, &current)
}

func ContainsToken(candidate string, query string) bool {
	candidateToken := Token(candidate)
	queryToken := Token(query)
	if candidateToken != "" && queryToken != "" && (candidateToken == queryToken || strings.Contains(candidateToken, queryToken)) {
		return true
	}
	return ContainsTerms(candidate, query)
}

func ContainsTerms(candidate string, query string) bool {
	candidateTerms := Terms(candidate)
	queryTerms := Terms(query)
	if len(candidateTerms) == 0 || len(queryTerms) == 0 {
		return false
	}
	for _, queryTerm := range queryTerms {
		if !termContained(candidateTerms, queryTerm) {
			return false
		}
	}
	return true
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
	return uniquePartialMatch(query, candidates)
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

func uniquePartialMatch(query string, candidates []Candidate) string {
	matched := ""
	for _, candidate := range candidates {
		if !candidateMatches(query, candidate) {
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

func candidateMatches(query string, candidate Candidate) bool {
	for _, value := range []string{candidate.Key, candidate.Name} {
		if ContainsToken(value, query) {
			return true
		}
	}
	return false
}

func appendSearchTerm(terms []string, current *strings.Builder) []string {
	if current.Len() == 0 {
		return terms
	}
	terms = append(terms, current.String())
	current.Reset()
	return terms
}

func termContained(candidateTerms []string, queryTerm string) bool {
	for _, candidateTerm := range candidateTerms {
		if strings.Contains(candidateTerm, queryTerm) {
			return true
		}
	}
	return false
}

func isSearchRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}
