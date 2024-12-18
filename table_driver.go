package mysql

type TypeOption string

type DBMSExternalDriver interface {
	EncodeValueByTypeName(name string, src any, buf []byte) ([]byte, error)
	DecodeValueByTypeName(name string, src []byte) (any, error)
	ScanValueByTypeName(name string, src []byte, dest any) error
}
