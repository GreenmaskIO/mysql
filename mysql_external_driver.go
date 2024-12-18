package mysql

import (
	"database/sql/driver"
	"strconv"
	"time"
)

const (
	// Numeric types
	TypeTinyInt   = "tinyint"
	TypeSmallInt  = "smallint"
	TypeMediumInt = "mediumint"
	TypeInt       = "int"
	TypeBigInt    = "bigint"
	TypeDecimal   = "decimal"
	TypeNumeric   = "numeric"
	TypeFloat     = "float"
	TypeDouble    = "double"
	TypeReal      = "real"
	TypeBit       = "bit"

	// Date and time types
	TypeDate      = "date"
	TypeDateTime  = "datetime"
	TypeTimestamp = "timestamp"
	TypeTime      = "time"
	TypeYear      = "year"

	// String types
	TypeChar    = "char"
	TypeVarChar = "varchar"

	// Text types
	TypeTinyText   = "tinytext"
	TypeText       = "text"
	TypeMediumText = "mediumtext"
	TypeLongText   = "longtext"

	// Binary types
	TypeBinary    = "binary"
	TypeVarBinary = "varbinary"

	// Blob types
	TypeTinyBlob   = "tinyblob"
	TypeBlob       = "blob"
	TypeMediumBlob = "mediumblob"
	TypeLongBlob   = "longblob"

	// Special string types
	TypeEnum = "enum"
	TypeSet  = "set"

	// Spatial types
	TypeGeometry           = "geometry"
	TypePoint              = "point"
	TypeLineString         = "linestring"
	TypePolygon            = "polygon"
	TypeMultiPoint         = "multipoint"
	TypeMultiLineString    = "multilinestring"
	TypeMultiPolygon       = "multipolygon"
	TypeGeometryCollection = "geometrycollection"

	// JSON type
	TypeJSON = "json"
)

var (
	// TypeOptionUnsigned - unsigned type
	TypeOptionUnsigned TypeOption = "unsigned"
)

type ExternalDriver struct {
	// con - dummy connection to mysql that used the internal to encode values
	con *mysqlConn
	Loc *time.Location
}

func NewExternalDriver() *ExternalDriver {
	return &ExternalDriver{
		con: &mysqlConn{
			buf:              newBuffer(nil),
			maxAllowedPacket: maxPacketSize,
			cfg: &Config{
				InterpolateParams: true,
			},
		},
		Loc: time.Now().Location(),
	}
}

func (e *ExternalDriver) WithLocation(loc *time.Location) *ExternalDriver {
	e.Loc = loc
	return e
}

func (e *ExternalDriver) EncodeValueByTypeName(name string, src any, buf []byte, opts ...TypeOption) ([]byte, error) {
	res, err := e.con.interpolateParams("?", []driver.Value{src})
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		res = res[1 : len(res)-1]
	}
	return []byte(res), nil
}

func (e *ExternalDriver) DecodeValueByTypeName(name string, src []byte, opts ...TypeOption) (any, error) {
	switch name {
	case TypeTimestamp,
		TypeDateTime,
		TypeDate:
		return parseDateTime(src, e.Loc)

	case TypeTinyInt, TypeSmallInt, TypeMediumInt, TypeInt, TypeBigInt, TypeTime, TypeYear:
		//case fieldTypeTiny, fieldTypeShort, fieldTypeInt24, fieldTypeYear, fieldTypeLong:
		if hasOptions(opts, TypeOptionUnsigned) {
			return strconv.ParseUint(string(src), 10, 64)
		} else {
			return strconv.ParseInt(string(src), 10, 64)
		}

	case TypeFloat, TypeDouble, TypeReal:
		return strconv.ParseFloat(string(src), 64)
	case TypeChar, TypeVarChar, TypeTinyText, TypeText, TypeMediumText, TypeLongText:
		return string(src), nil
	}
	return src, nil
}

func (e *ExternalDriver) ScanValueByTypeName(name string, src []byte, dest any) error {

	//TODO implement me
	panic("implement me")
}

func hasOptions(opts []TypeOption, opt TypeOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false

}
