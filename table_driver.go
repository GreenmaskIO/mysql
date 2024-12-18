package mysql

type TypeOption string

type DBMSExternalDriver interface {
	EncodeValueByTypeName(name string, src any, buf []byte, opts ...TypeOption) ([]byte, error)
	DecodeValueByTypeName(name string, src []byte, opts ...TypeOption) (any, error)
	ScanValueByTypeName(name string, src []byte, dest any, opts ...TypeOption) error
}
