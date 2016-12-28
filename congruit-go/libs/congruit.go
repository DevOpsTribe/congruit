package congruit

type Work struct {
	Name    string
	Command string
}

type Place struct {
	Name    string
	Command string
}

type WorkPlace struct {
	Name   string
	Works  []string
	Places []string
}

func (w *Work) DoWork() {
	println("[CONG] doing " + w.Name)

}
