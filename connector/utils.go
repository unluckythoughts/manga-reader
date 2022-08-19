package connector

type strFloat string

func (sf *strFloat) UnmarshalJSON(data []byte) error {
	*sf = strFloat(data)

	return nil
}
