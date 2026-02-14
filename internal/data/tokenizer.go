package data

type Tokenizer interface {
	Decode(data []byte) (string, error)
	Encode(data string) ([]byte, error)
}
