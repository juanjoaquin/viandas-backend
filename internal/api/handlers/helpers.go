package handlers

func optionalString(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
