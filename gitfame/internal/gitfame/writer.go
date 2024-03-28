package gitfame

type Writer interface {
	Write(statistics *[]Statistics)
}
