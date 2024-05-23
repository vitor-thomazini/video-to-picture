package domain

type StorageBatch struct {
	LatestIndex int
	LatestText  string
	Data        map[string][]Resource
}

func NewStorageBatch() StorageBatch {
	return StorageBatch{
		LatestIndex: 0,
		LatestText:  "",
		Data:        make(map[string][]Resource),
	}
}

func (l StorageBatch) GetData() map[string][]Resource {
	if len(l.Data) <= 0 {
		return make(map[string][]Resource)
	}
	return l.Data
}
