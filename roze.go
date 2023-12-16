package rozer

type Roze struct {
	data  []any
	names map[string]int
}

func New() *Roze {
	return &Roze{}
}

func (r *Roze) Len() int {
	return len(r.data)
}

func (r *Roze) Append(value any) *Roze {
	r.data = append(r.data, value)
	return r
}

func (r *Roze) Put(name string, value any) *Roze {
	if r.names == nil {
		r.names = make(map[string]int)
	}
	p, ok := r.names[name]
	if ok {
		r.data[p] = value
		return r
	}
	r.names[name] = len(r.data)
	r.data = append(r.data, value)
	return r
}

func (r *Roze) Get(name string) any {
	return r.data[r.names[name]]
}

func (r *Roze) At(i int) any {
	if i >= len(r.data) {
		return nil
	}
	if i < 0 {
		i = len(r.data) + i
	}
	return r.data[i]
}

func (r *Roze) Set(i int, value any) *Roze {
	if i >= len(r.data) {
		r.data = append(r.data, make([]any, i-len(r.data)+1)...)
	}
	r.data[i] = value
	return r
}
