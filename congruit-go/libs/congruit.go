package congruit

type Work struct {
	Name        string
	DoAfter     string
	Command     string
	Idempotency string
}

type Place struct {
	Name    string
	Command string
}

func (w *Work) DoWork() {
	println("[CONG] doing " + w.Name)
}
