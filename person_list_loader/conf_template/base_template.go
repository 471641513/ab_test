package conf_template

type ListTemplate interface {
	ParseConfig(configStr string) (confMap map[int64]map[int64]bool, err error)
	ParseList(listStr string) (err error)
	Merge(IFace interface{}) (err error)
}
