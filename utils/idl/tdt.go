package idl

type TypeDefinitionTable struct {
	Types   [][]byte
	Indexes map[string]int
}

func (tdt *TypeDefinitionTable) Add(t Type, bs []byte) {
	if tdt.Has(t) {
		return 
	}
	i := len(tdt.Types)
	tdt.Indexes[t.String()] = i
	tdt.Types = append(tdt.Types, bs)
}

func (tdt *TypeDefinitionTable) Has(t Type) bool{
	_, ok := tdt.Indexes[t.String()]
	return ok
}