package need

import "sync"

type NeedData struct {
	data interface{}
	once sync.Once
	wg   sync.WaitGroup
}

func (n *NeedData) Need(retrieve func() interface{}) func() interface{} {
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
