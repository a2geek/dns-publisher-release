package triggers

import "time"

func newRefreshTrigger(refresh string) (Trigger, error) {
	duration, err := time.ParseDuration(refresh)
	if err != nil {
		return nil, err
	}
	return &refreshTrigger{
		duration: duration,
	}, nil
}

type refreshTrigger struct {
	duration time.Duration
}
type refreshTick struct {
	data time.Time
}

func (t *refreshTrigger) Start() (<-chan interface{}, error) {
	ch := make(chan interface{})
	go func() {
		for t := range time.Tick(t.duration) {
			ch <- refreshTick{
				data: t,
			}
		}
	}()
	return ch, nil
}
