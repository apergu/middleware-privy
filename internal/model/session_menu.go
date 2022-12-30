package model

type SessionMenu map[string]int8

func (s SessionMenu) FindAccess(menu string) int8 {
	access, ok := s[menu]
	if !ok {
		return 0
	}

	return access
}

func (s SessionMenu) IsHasMenu(menu string) bool {
	return s.FindAccess(menu) > 0
}
