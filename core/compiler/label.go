package compiler

type Label struct {
	address    int
	patchWhere int
}

type LoopContext struct {
	breakLabel    *Label
	continueLabel *Label
}

func (c *Compiler) newLabel() *Label {
	return &Label{
		address:    -1,
		patchWhere: -1,
	}
}

func (c *Compiler) markLabel(label *Label) {
	label.address = (len(c.bc.Code))

	if label.patchWhere != -1 {
		c.backpatch(label)
	}
}

func (c *Compiler) backpatch(label *Label) {
	if label.address < 0 || label.address > 0xFFFF {
		panic("backpatch address does not fit in uint16")
	}

	addr := uint16(label.address)
	c.bc.Code[label.patchWhere] = byte(addr >> 8)
	c.bc.Code[label.patchWhere+1] = byte(addr)
}
