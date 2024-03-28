//go:build !solution

package tparallel

type T struct {
	finished   chan struct{}
	bar        chan struct{}
	parent     *T
	subtest    []*T
	isParallel bool
}

func NewT(t *T) *T {
	return &T{
		finished: make(chan struct{}),
		bar:      make(chan struct{}),
		parent:   t,
		subtest:  make([]*T, 0),
	}
}

func TMake() *T {
	return NewT(nil)
}

func (t *T) Parallel() {
	t.isParallel = true
	t.parent.subtest = append(t.parent.subtest, t)

	t.finished <- struct{}{}
	<-t.parent.bar
}

func (t *T) tRunner(subtest func(t *T)) {
	subtest(t)

	if len(t.subtest) > 0 {
		close(t.bar)

		for _, sub := range t.subtest {
			<-sub.finished
		}
	}

	if t.isParallel {
		t.parent.finished <- struct{}{}
	}

	t.finished <- struct{}{}
}

func (t *T) Run(subtest func(t *T)) {
	st := NewT(t)
	go st.tRunner(subtest)
	<-st.finished
}

func Run(topTests []func(t *T)) {
	main := TMake()

	for _, grow := range topTests {
		main.Run(grow)
	}

	close(main.bar)

	if len(main.subtest) > 0 {
		<-main.finished
	}
}
