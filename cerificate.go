package agent

import (
	"encoding/hex"

	"github.com/fxamacker/cbor/v2"
)

const (
	Empty uint64 = iota
	Fork
	Labeled
	Leaf
	Pruned
)

func LookUp(paths [][]byte, cert []byte)([]byte, error) {
	var certificate interface{}
	err := cbor.Unmarshal(cert, &certificate)
	if err != nil{
		return nil,err
	}
	certi_struct, ok := certificate.(map[interface{}]interface{})
	if !ok{
		return nil, nil
	}
	for k, v := range(certi_struct){
		k_value, ok := k.(string)
		if ok && k_value == "tree"{
			tree_value, ok := v.([]interface{})
			if !ok {
				continue
			}
			
			return lookupPath(paths, tree_value)

		}
	}
	
	return nil, nil
}

func lookupPath(paths [][]byte, tree []interface{})([]byte, error) {
	if tree == nil{
		return nil, nil
	}
	offset := 0
	if len(paths) == 0{
		if tree == nil{
			return nil, nil
		}
		tree_0, ok := tree[0].(uint64)
		if !ok{
			return nil, nil
		}
		if tree_0 == Leaf{
			tree_1, ok := tree[1].([]byte)
			if !ok{
				return nil, nil
			}
			return tree_1, nil
		} else {
			return nil, nil
		}
	}
	label := paths[0]


	t_flatten, _ := flattenForks(tree)
	t, _ := findLabel(label, t_flatten)
	offset += 1
	return lookupPath(paths[offset:], t)
}

func flattenForks(tree []interface{}) ([][]interface{}, error){
	tree_0, ok := tree[0].(uint64)
	if !ok{
		return nil, nil
	}
	if tree_0 == Empty{
		return [][]interface{}{}, nil
	} else if tree_0 == Fork{
		t_1, ok := tree[1].([]interface{})
		if !ok{
			return nil, nil
		}
		t_2, ok := tree[2].([]interface{})
		if !ok{
			return nil, nil
		}
		left, _ := flattenForks(t_1)
		right, _ := flattenForks(t_2)
		return append(left, right...), nil
	} else {
		return [][]interface{}{tree}, nil
	}
}

func findLabel(label []byte, trees [][]interface{}) ([]interface{}, error) {
	if len(trees) == 0{
		return nil, nil
	}

	for _, t := range(trees) {
		t_0, ok := t[0].(uint64)
		if !ok{
			return nil, nil
		}
		if t_0 == Labeled{
			t_1, ok := t[1].([]byte)
			
			if !ok {
				
				return nil, nil
			}
			if (hex.EncodeToString(t_1) != hex.EncodeToString(label)) {
				//fmt.Println(label, "   error")
				continue
			}
			t_2, ok := t[2].([]interface{})
			if !ok{
				return nil, nil
			}
			return t_2, nil
		}
	}
	return nil, nil
}


type Certificate struct {
	Tree       Tree   `cbor:"tree"`
	Signature  []byte `cbor:"signature"`
	Delegation []byte `cbor:"delegation"`
}

type Tree struct {
	_   struct{} `cbor:",toarray"`
	sym byte
	a   []byte
	b   []byte
}
