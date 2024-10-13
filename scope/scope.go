package scope

const (
	S_ALL uint8 = 1 << iota
	S_ACCESS
	S_CREATE
	S_DELETE
	S_GET
	S_LIST
	S_MODIFY
)

func ConfHasOpt(conf byte, compare byte) bool {
	return conf&compare != 0
}

func CreatePermission(conf byte) Permission {
	data := Permission{
		All:    ConfHasOpt(conf, S_ALL),
		Access: ConfHasOpt(conf, S_ACCESS),
		Create: ConfHasOpt(conf, S_CREATE),
		Delete: ConfHasOpt(conf, S_DELETE),
		Get:    ConfHasOpt(conf, S_GET),
		List:   ConfHasOpt(conf, S_LIST),
		Modify: ConfHasOpt(conf, S_MODIFY),
	}
	return data
}
