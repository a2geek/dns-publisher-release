package processors

type Processor interface {
	Run(actionChan chan<- Action)
}

type Action interface {
	Name() string
	Act()
}
