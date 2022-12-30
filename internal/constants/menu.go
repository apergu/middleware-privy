package constants

const (
	MenuCreate int8 = 1 << iota
	MenuRead
	MenuUpdate
	MenuDelete
)

func IsAllowCreate(access int8) bool {
	return FindStatus(access, MenuCreate)
}

func IsAllowRead(access int8) bool {
	return FindStatus(access, MenuRead)
}

func IsAllowUpdate(access int8) bool {
	return FindStatus(access, MenuUpdate)
}

func IsAllowDelete(access int8) bool {
	return FindStatus(access, MenuDelete)
}
