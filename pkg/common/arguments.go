package common

type Arguments map[string]string

func (t *Arguments) SetStrArgument(name, value string) {
	(*t)[name] = value
}

func (t *Arguments) GetStrArgument(name string) (string, bool) {
	if !t.HasKey(name) {
		return "", false
	}
	value := (*t)[name]
	return value, true
}

func (t *Arguments) HasKey(name string) bool {
	if t == nil {
		return false
	}
	_, ok := (*t)[name]
	return ok
}
