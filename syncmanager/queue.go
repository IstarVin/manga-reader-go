package syncmanager

type Queue struct {
	Index    int
	Function func(...any)
	Args     []any
}

type QueueList struct {
	List []*Queue
}

func (l *QueueList) FindQueue(index int) *Queue {
	for _, q := range l.List {
		if q.Index == index {
			return q
		}
	}
	return nil
}
