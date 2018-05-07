package need

import "sync"

type NeedData struct {
	data interface{}
	once sync.Once
	wg   sync.WaitGroup
}

type Blah func() interface{}

func (n *NeedData) Need(retrieve Blah) Blah {
	n.once.Do(func() {
		n.wg.Add(1)
		go func() {
			n.data = retrieve()
			n.wg.Done()
		}()
	})

	return func() interface{} {
		n.wg.Wait()
		return n.data
	}
}
