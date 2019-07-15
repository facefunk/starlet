package scarlet

func filterTags(rules []*CSSRule, tags []string) []*CSSRule {
	if len(tags) == 0 {
		return rules
	}
	var out []*CSSRule
	for _, rule := range rules {
		rule.Duplicates = filterTags(rule.Duplicates, tags)
		include := true
	loop:
		for _, part := range rule.Selector {
			if part.Type != ElementSelector {
				continue
			}
			include = false
			for _, t := range tags {
				if t == part.Name {
					include = true
					break loop
				}
			}
		}
		if include {
			out = append(out, rule)
			continue
		}
		if len(rule.Duplicates) > 0 {
			r := rule.Duplicates[0]
			r.Duplicates = append(r.Duplicates, rule.Duplicates[1:]...)
			out = append(out, r)
		}
	}
	return out
}
