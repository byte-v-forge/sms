export function normalizeDisplayName(value: string) {
  return value.replace(/\s*([,/|;])\s*/g, '$1 ').replace(/\s+/g, ' ').trim();
}

export function normalizeChoiceToken(value: string) {
  return value.trim().toLowerCase().replace(/[^a-z0-9]+/g, '');
}

export function searchTerms(value: string) {
  return value.toLowerCase().match(/[a-z0-9]+/g) || [];
}

export function termsContain(candidateTerms: string[], queryTerms: string[]) {
  return queryTerms.every((query) => candidateTerms.some((candidate) => candidate.includes(query)));
}

export function titleSlug(value: string) {
  return value.split(/[-_.\s]+/).filter(Boolean).map(titlePart).join(' ');
}

function titlePart(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1).toLowerCase();
}
