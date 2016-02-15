package block

/*

0-3  # records





*/

type Block struct {
	RecordTable []int
	Records     []byte
}

type Field struct {
	FieldTable []int
	Fields     []byte
}
