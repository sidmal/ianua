package security

type None struct{}

func newNoneSigner() Signer {
	return &None{}
}

func (m *None) GetSignature(_ string, _ map[string]interface{}) (string, error) {
	return "", nil
}
