package properties

type CanValidateValue interface {
	ValidateICalValue() error
}

type CanDecodeValue interface {
	DecodeICalValue(value string) error
}

type CanDecodeParams interface {
	DecodeICalParams(value map[ParameterName]string) error
}

type CanEncodeTag interface {
	EncodeICalTag() (string, error)
}

type CanEncodeValue interface {
	EncodeICalValue() (string, error)
}

type CanEncodeName interface {
	EncodeICalName() (PropertyName, error)
}

type CanEncodeParams interface {
	EncodeICalParams() (map[ParameterName]string, error)
}
