package vector

import "github.com/koykov/bytealg"

type Type int

const (
	TypeUnk Type = iota
	TypeNull
	TypeObj
	TypeArr
	TypeStr
	TypeNum
	TypeBool
	TypeAttr
)

type Node struct {
	typ      Type
	depth    int
	vecPtr   uintptr
	key, val bytealg.Byteptr
	off, len int
}
