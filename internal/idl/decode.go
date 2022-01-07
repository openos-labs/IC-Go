package idl

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/aviate-labs/leb128"
)


type typePair struct {
	Type int64
	single_value int64
	pair_value [][]big.Int
}

func _Decode(raw_table []typePair, index int64) (Type, error){
	pair := raw_table[index]
	switch pair.Type {
	case optType:
		tid := pair.single_value
		var v Type
		var err error
		if tid >= 0 {
			if int(tid) >= len(raw_table) {
				return nil, nil
			}
			v, err = _Decode(raw_table, tid)
		} else{
			v, err = getType(tid)
		}
		if err != nil {
			return nil, err
		}
		return &Opt{v}, nil
	case vecType:
		tid := pair.single_value
		var v Type
		var err error
		if tid >= 0 {
			if int(tid) >= len(raw_table) {
				return nil, nil
			}
			v, err = _Decode(raw_table, tid)
		} else{
			v, err = getType(tid)
		}
		if err != nil {
			return nil, err
		}
		return &Vec{v}, nil
		//tds = append(tds, &Vec{v})
	case recType:
		l := len(pair.pair_value)
		var fields []Field
		for i := 0; i < int(l); i++ {
			h := pair.pair_value[i][0]
			tid := pair.pair_value[i][1]
			var v Type
			var err error
			if tid.Int64() >= 0 {
				if int(tid.Int64()) >= len(raw_table) {
					return nil, nil
				}
				v, err = _Decode(raw_table, tid.Int64())
			} else{
				v, err = getType(tid.Int64())
			}
			
			if err != nil {
				return nil, err
			}
			fields = append(fields, Field{
				Name: h.String(),
				Type: v,
			})
		}
		return &Rec{Fields: fields}, nil
	case varType:
		l := len(pair.pair_value)
		var fields []Field
		for i := 0; i < int(l); i++ {
			h := pair.pair_value[i][0]
			tid := pair.pair_value[i][1]
			var v Type
			var err error
			if tid.Int64() >= 0 {
				if int(tid.Int64()) >= len(raw_table) {
					return nil, nil
				}
				v, err = _Decode(raw_table, tid.Int64())
			} else{
				v, err = getType(tid.Int64())
			}
			
			if err != nil {
				return nil, err
			}
			fields = append(fields, Field{
				Name: h.String(),
				Type: v,
			})
		}
		return &Variant{Fields: fields}, nil
	}
	return nil, nil
}

func Decode(bs []byte) ([]Type, []interface{}, error) {
	fmt.Println("aaa", len(bs))
	fmt.Println(hex.EncodeToString(bs))
	if len(bs) == 0 {
		return nil, nil, &FormatError{
			Description: "empty",
		}
	}

	r := bytes.NewReader(bs)

	{ // 'DIDL'

		magic := make([]byte, 4)
		n, err := r.Read(magic)
		if err != nil {
			return nil, nil, err
		}
		if n < 4 {
			return nil, nil, &FormatError{
				Description: "no magic bytes",
			}
		}
		if !bytes.Equal(magic, []byte{'D', 'I', 'D', 'L'}) {
			return nil, nil, &FormatError{
				Description: "wrong magic bytes",
			}
		}
	}
	var raw_table []typePair
	{ // T
		tdtl, err := leb128.DecodeUnsigned(r)
		if err != nil {
			return nil, nil, err
		}
		for i := 0; i < int(tdtl.Int64()); i++ {
			tid, err := leb128.DecodeSigned(r)
			if err != nil {
				return nil, nil, err
			}
			switch tid.Int64() {
			case optType:
				t, err := leb128.DecodeSigned(r)
				if err != nil {
					return nil, nil, err
				}
				raw_table = append(raw_table, typePair{
					Type: tid.Int64(),
					single_value: t.Int64(),
					pair_value: nil,
				})
			case vecType:
				t, err := leb128.DecodeSigned(r)
				if err != nil {
					return nil, nil, err
				}
				raw_table = append(raw_table, typePair{
					Type: tid.Int64(),
					single_value: t.Int64(),
					pair_value: nil,
				})
			case recType:
				l, err := leb128.DecodeUnsigned(r)
				if err != nil {
					return nil, nil, err
				}
				var fields [][]big.Int
				for i := 0; i < int(l.Int64()); i++ {
					h, err := leb128.DecodeUnsigned(r)
					if err != nil {
						return nil, nil, err
					}
					t, err := leb128.DecodeSigned(r)
					if err != nil {
						return nil, nil, err
					}
					fields = append(fields, []big.Int{*h, *t})
				}
				raw_table = append(raw_table, typePair{
					Type: tid.Int64(),
					single_value: 99,
					pair_value: fields,
				})	
			case varType:	
				l, err := leb128.DecodeUnsigned(r)
				if err != nil {
					return nil, nil, err
				}
				var fields [][]big.Int
				for i := 0; i < int(l.Int64()); i++ {
					h, err := leb128.DecodeUnsigned(r)
					if err != nil {
						return nil, nil, err
					}
					t, err := leb128.DecodeSigned(r)
					if err != nil {
						return nil, nil, err
					}
					fields = append(fields, []big.Int{*h, *t})
				}
				raw_table = append(raw_table, typePair{
					Type: tid.Int64(),
					single_value: 99,
					pair_value: fields,
				})	
			case funcType:
				// la, err := leb128.DecodeUnsigned(r)
				// if err != nil {
				// 	return nil, nil, err
				// }
				// var args []Type
				// for i := 0; i < int(la.Int64()); i++ {
				// 	tid, err = leb128.DecodeSigned(r)
				// 	if err != nil {
				// 		return nil, nil, err
				// 	}
				// 	v, err := getType(tid.Int64(), tds)
				// 	if err != nil {
				// 		return nil, nil, err
				// 	}
				// 	args = append(args, v)
				// }
				// lr, err := leb128.DecodeUnsigned(r)
				// if err != nil {
				// 	return nil, nil, err
				// }
				// var rets []Type
				// for i := 0; i < int(lr.Int64()); i++ {
				// 	tid, err = leb128.DecodeSigned(r)
				// 	if err != nil {
				// 		return nil, nil, err
				// 	}
				// 	v, err := getType(tid.Int64(), tds)
				// 	if err != nil {
				// 		return nil, nil, err
				// 	}
				// 	rets = append(rets, v)
				// }
				// l, err := leb128.DecodeUnsigned(r)
				// if err != nil {
				// 	return nil, nil, err
				// }
				// ann := make([]byte, l.Int64())
				// if _, err := r.Read(ann); err != nil {
				// 	return nil, nil, err
				// }
				// var anns []string
				// if len(ann) != 0 {
				// 	anns = append(anns, string(ann))
				// }
				// tds = append(tds, &Func{
				// 	ArgTypes:    args,
				// 	RetTypes:    rets,
				// 	Annotations: anns,
				// })
			case serviceType:
				// l, err := leb128.DecodeUnsigned(r)
				// if err != nil {
				// 	return nil, nil, err
				// }
				// var methods []Method
				// for i := 0; i < int(l.Int64()); i++ {
				// 	lm, err := leb128.DecodeUnsigned(r)
				// 	if err != nil {
				// 		return nil, nil, err
				// 	}
				// 	name := make([]byte, lm.Int64())
				// 	n, err := r.Read(name)
				// 	if err != nil {
				// 		return nil, nil, err
				// 	}
				// 	if n != int(lm.Int64()) {
				// 		return nil, nil, fmt.Errorf("invalid method name: %d", bs)
				// 	}

				// 	tid, err = leb128.DecodeSigned(r)
				// 	if err != nil {
				// 		return nil, nil, err
				// 	}
				// 	v, err := getType(tid.Int64(), tds)
				// 	if err != nil {
				// 		return nil, nil, err
				// 	}
				// 	f, ok := v.(*Func)
				// 	if !ok {
				// 		fmt.Println(reflect.TypeOf(v))
				// 	}
				// 	methods = append(methods, Method{
				// 		Name: string(name),
				// 		Func: f,
				// 	})
				// }
				// tds = append(tds, &Service{
				// 	methods: methods,
				// })
			}
		}
	}
	var tds []Type
	{ // T
		for i, _ := range(raw_table){
			t, err := _Decode(raw_table, int64(i))
			if err != nil {
				return nil, nil, err
			}
			tds = append(tds, t)
		}
	}
	// var tds []Type
	// { // T
	// 	tdtl, err := leb128.DecodeUnsigned(r)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	for i := 0; i < int(tdtl.Int64()); i++ {
	// 		tid, err := leb128.DecodeSigned(r)
	// 		if err != nil {
	// 			return nil, nil, err
	// 		}
	// 		switch tid.Int64() {
	// 		case optType:
	// 			tid, err := leb128.DecodeSigned(r)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			v, err := getType(tid.Int64(), tds)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			tds = append(tds, &Opt{v})
	// 		case vecType:
	// 			tid, err := leb128.DecodeSigned(r)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			v, err := getType(tid.Int64(), tds)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			tds = append(tds, &Vec{v})
	// 		case recType:
	// 			l, err := leb128.DecodeUnsigned(r)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			var fields []Field
	// 			for i := 0; i < int(l.Int64()); i++ {
	// 				h, err := leb128.DecodeUnsigned(r)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				tid, err := leb128.DecodeSigned(r)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				v, err := getType(tid.Int64(), tds)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				fields = append(fields, Field{
	// 					Name: h.String(),
	// 					Type: v,
	// 				})
	// 			}
	// 			tds = append(tds, &Rec{Fields: fields})
	// 		case varType:
	// 			l, err := leb128.DecodeUnsigned(r)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			var fields []Field
	// 			for i := 0; i < int(l.Int64()); i++ {
	// 				h, err := leb128.DecodeUnsigned(r)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				tid, err := leb128.DecodeSigned(r)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				v, err := getType(tid.Int64(), tds)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				fields = append(fields, Field{
	// 					Name: h.String(),
	// 					Type: v,
	// 				})
	// 			}
	// 			tds = append(tds, &Variant{Fields: fields})
	// 		case funcType:
	// 			la, err := leb128.DecodeUnsigned(r)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			var args []Type
	// 			for i := 0; i < int(la.Int64()); i++ {
	// 				tid, err = leb128.DecodeSigned(r)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				v, err := getType(tid.Int64(), tds)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				args = append(args, v)
	// 			}
	// 			lr, err := leb128.DecodeUnsigned(r)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			var rets []Type
	// 			for i := 0; i < int(lr.Int64()); i++ {
	// 				tid, err = leb128.DecodeSigned(r)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				v, err := getType(tid.Int64(), tds)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				rets = append(rets, v)
	// 			}
	// 			l, err := leb128.DecodeUnsigned(r)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			ann := make([]byte, l.Int64())
	// 			if _, err := r.Read(ann); err != nil {
	// 				return nil, nil, err
	// 			}
	// 			var anns []string
	// 			if len(ann) != 0 {
	// 				anns = append(anns, string(ann))
	// 			}
	// 			tds = append(tds, &Func{
	// 				ArgTypes:    args,
	// 				RetTypes:    rets,
	// 				Annotations: anns,
	// 			})
	// 		case serviceType:
	// 			l, err := leb128.DecodeUnsigned(r)
	// 			if err != nil {
	// 				return nil, nil, err
	// 			}
	// 			var methods []Method
	// 			for i := 0; i < int(l.Int64()); i++ {
	// 				lm, err := leb128.DecodeUnsigned(r)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				name := make([]byte, lm.Int64())
	// 				n, err := r.Read(name)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				if n != int(lm.Int64()) {
	// 					return nil, nil, fmt.Errorf("invalid method name: %d", bs)
	// 				}

	// 				tid, err = leb128.DecodeSigned(r)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				v, err := getType(tid.Int64(), tds)
	// 				if err != nil {
	// 					return nil, nil, err
	// 				}
	// 				f, ok := v.(*Func)
	// 				if !ok {
	// 					fmt.Println(reflect.TypeOf(v))
	// 				}
	// 				methods = append(methods, Method{
	// 					Name: string(name),
	// 					Func: f,
	// 				})
	// 			}
	// 			tds = append(tds, &Service{
	// 				methods: methods,
	// 			})
	// 		}
	// 	}
	// }

	tsl, err := leb128.DecodeUnsigned(r)
	if err != nil {
		return nil, nil, err
	}

	var ts []Type
	{ // I
		for i := 0; i < int(tsl.Int64()); i++ {
			tid, err := leb128.DecodeSigned(r)
			var t Type
			if err != nil {
				return nil, nil, err
			}
			if tid.Int64() < 0 {
				t, err = getType(tid.Int64())
				if err != nil {
					return nil, nil, err
				}
			} else {
				t = tds[int(tid.Int64())]
			}
			ts = append(ts, t)
		}
	}

	var vs []interface{}
	{ // M
		for i := 0; i < int(tsl.Int64()); i++ {
			v, err := ts[i].Decode(r)
			if err != nil {
				return nil, nil, err
			}
			vs = append(vs, v)
		}
	}

	if r.Len() != 0 {
		return nil, nil, fmt.Errorf("too long")
	}
	fmt.Println("vs   ", len(vs))
	return ts, vs, nil
}
