package opcode

type Opcode = int

const (
	LOAD_CONST = iota
	STORE_NAME
	LOAD_NAME
	CALL_FUNCTION
	RELATIVE_JUMP_IF_TRUE
	RELATIVE_JUMP
)
