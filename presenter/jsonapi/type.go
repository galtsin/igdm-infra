package jsonapi

type Type struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Resource struct {
	Type
	Attributes interface{} `json:"attributes,omitempty"`
}

type With struct {
	res [][]*Resource
	len int
}

func NewWith() *With {
	return &With{
		res: make([][]*Resource, 0, 3),
	}
}

func (w *With) append(res []*Resource) {
	w.res = append(w.res, res)
	w.len = w.len + len(res)
}

func (w *With) all() []*Resource {
	res := make([]*Resource, 0, w.len)
	for _, list := range w.res {
		res = append(res, list...)
	}

	return res
}
