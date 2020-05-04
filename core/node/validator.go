package node

type NullValidator struct {

}

func (v *NullValidator) Validate(key string, value []byte) error {
	return nil
}

func (v *NullValidator) Select(key string, values [][]byte) (int, error) {
	return 0, nil
}
