package syncmanager

var SyncManager *Concurrency

func Init(concurrent int) {
	SyncManager = &Concurrency{}
	SyncManager.Init(concurrent)
}

type Concurrency struct {
	State    int
	Multiple int
	QueueNo  int

	QueueList *QueueList

	wait      bool
	ch        chan string
	finish    chan string
	currIndex int
}

func (c *Concurrency) Init(concurrent int) {
	c.State = IDLE
	c.Multiple = concurrent
	c.QueueList = &QueueList{}
	c.QueueList.List = []*Queue{}
	c.ch = make(chan string)
	go SyncManager.daemon()
}

func (c *Concurrency) AddQueue(f func(...any), args ...any) {
	c.QueueList.List = append(c.QueueList.List, &Queue{c.currIndex, f, args})
	c.currIndex++
	c.ch <- ADDQUEUE
}

func (c *Concurrency) daemon() {
	daemonCH := make(chan string)

	var index, queueNo int

	for {
		select {
		case <-daemonCH:
			queueNo--
		case <-c.ch:
		}

		if len(c.QueueList.List) > 0 && queueNo <= c.Multiple && len(c.QueueList.List) > index {
			go func(i int) {
				queue := c.QueueList.FindQueue(i)
				queue.Function(queue.Args...)
				queue.Function = nil
				daemonCH <- DONE
			}(index)
			index++
			queueNo++
		}

		if c.checkIfAllIsDone() && c.wait {
			c.State = IDLE
			c.finish <- DONE
			break
		}
	}
}

func (c *Concurrency) WaitFinish() {
	c.finish = make(chan string)
	c.wait = true
	<-c.finish
}

func (c *Concurrency) checkIfAllIsDone() bool {
	for _, queue := range c.QueueList.List {
		if queue.Function != nil {
			return false
		}
	}
	return true
}
