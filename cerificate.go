package agent

import (
	"bytes"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

const (
	Empty byte = iota
	Fork
	Labeled
	Leaf
	Pruned
)

func LookUp(paths [][]byte, cert []byte)([]byte, error) {
	//certificate := new(Certificate)
	var certificate interface{}
	err := cbor.Unmarshal(cert, &certificate)
	if err != nil{
		return nil,err
	}
	certi_struct, ok := certificate.(map[string]interface{})
	if ok {
		fmt.Println(certi_struct["tree"])
	}
	return nil, nil//lookupPath(paths,&certificate.Tree)
}

func lookupPath(paths [][]byte, tree *Tree) ([]byte, error) {
	offset := 0
	if len(paths) == 0 {
		//todo:不确定大端小端到时候看看吧
		fmt.Printf("[certificate]:the tree label d is %d", tree.sym)
		if tree.sym == Leaf {
			return tree.a,nil
		} else {
			return nil,fmt.Errorf("can not find the path %x",paths)
		}
	}
	trees, err := flattenForks(tree)
	if err != nil {
		return nil, err
	}
	t,err := findLabel(paths[0],trees)
	if t != nil {
		offset++
		return lookupPath(paths[offset:],t)
	}
	return nil,fmt.Errorf("can not find the path %x",paths)
}

func flattenForks(tree *Tree) ([]*Tree, error) {
	var trees []*Tree
	if tree.sym == Fork {
		left := new(Tree)
		err := cbor.Unmarshal(tree.a, left)
		if err != nil {
			return trees, err
		}
		right := new(Tree)
		err = cbor.Unmarshal(tree.a, right)
		if err != nil {
			return trees, err
		}
		leftSubTree, err := flattenForks(left)
		if err != nil {
			return trees, err
		}
		rightSubTree, err := flattenForks(right)
		if err != nil {
			return trees, err
		}
		trees = append(trees, leftSubTree...)
		trees = append(trees, rightSubTree...)
	}
	return trees, nil
}

func findLabel(label []byte, trees []*Tree) (*Tree, error) {
	if len(trees) == 0 {
		return nil, nil
	}
	for _, t := range trees {
		if t.sym == Labeled {
			if bytes.Equal(label, t.a) {
				subTree := new(Tree)
				err := cbor.Unmarshal(t.b, subTree)
				if err != nil {
					return nil, err
				}
				return subTree, nil
			}
		}
	}
	return nil, fmt.Errorf("can not find lable : %x", string(label))
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
