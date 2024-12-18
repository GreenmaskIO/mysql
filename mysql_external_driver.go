package mysql

import (
	"fmt"
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

	TypeBoolean = "boolean"

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
	Loc *time.Location
}

func NewExternalDriver() *ExternalDriver {
	return &ExternalDriver{
		Loc: time.Now().Location(),
	}
}

func (e *ExternalDriver) WithLocation(loc *time.Location) *ExternalDriver {
	e.Loc = loc
	return e
}

func (e *ExternalDriver) EncodeValueByTypeName(name string, src any, buf []byte) ([]byte, error) {
	switch name {
	case TypeJSON:
		return encodeJson(src, buf)
	case TypeTimestamp,
		TypeDateTime,
		TypeDate:
		return encodeTimestamp(src, buf, e.Loc)
	case TypeTinyInt, TypeSmallInt, TypeMediumInt, TypeInt, TypeBigInt, TypeTime, TypeYear:
		return encodeInt64(src, buf)
	case TypeFloat, TypeDouble, TypeReal:
		return encodeFloat(src, buf)
	case TypeChar, TypeVarChar, TypeTinyText, TypeText, TypeMediumText, TypeLongText:
		return encodeString(src, buf)
	case TypeBoolean:
		return encodeBool(src, buf)
	case TypeBinary, TypeVarBinary, TypeTinyBlob, TypeBlob, TypeMediumBlob, TypeLongBlob:
		return encodeBinary(src, buf)
	case TypeEnum, TypeSet:
		return encodeEnum(src, buf)
	case TypeGeometry, TypePoint, TypeLineString, TypePolygon, TypeMultiPoint, TypeMultiLineString,
		TypeMultiPolygon, TypeGeometryCollection:
		return encodeGeometry(src, buf)
	case TypeDecimal, TypeNumeric:
		return encodeDecimal(src, buf)
	}
	return nil, fmt.Errorf("unsupported type %s", name)
}

func (e *ExternalDriver) DecodeValueByTypeName(name string, src []byte) (any, error) {
	// Consider opts pattern usage
	switch name {
	case TypeJSON:
		return string(src), nil
	case TypeTimestamp,
		TypeDateTime,
		TypeDate:
		return parseDateTime(src, e.Loc)
	case TypeTinyInt, TypeSmallInt, TypeMediumInt, TypeInt, TypeBigInt, TypeTime, TypeYear:
		// Here may be unsigned type consider to add it but it is likely redundant
		return strconv.ParseInt(string(src), 10, 64)
	case TypeFloat, TypeDouble, TypeReal:
		return strconv.ParseFloat(string(src), 64)
	case TypeChar, TypeVarChar, TypeTinyText, TypeText, TypeMediumText, TypeLongText:
		return string(src), nil
	case TypeBoolean:
		return decodeBool(src)
	case TypeBinary, TypeVarBinary, TypeTinyBlob, TypeBlob, TypeMediumBlob, TypeLongBlob:
		// I suspect there might be some hex encoding
		return src, nil
	case TypeEnum, TypeSet:
		return decodeEnum(src)
	case TypeGeometry, TypePoint, TypeLineString, TypePolygon, TypeMultiPoint, TypeMultiLineString,
		TypeMultiPolygon, TypeGeometryCollection:
		return decodeGeometry(src)
	case TypeDecimal, TypeNumeric:
		return decodeDecimal(src)
	}
	return nil, fmt.Errorf("unsupported type %s", name)
}

func (e *ExternalDriver) ScanValueByTypeName(name string, src []byte, dest any) error {
	switch name {
	case TypeJSON:
		return scanJson(src, dest)
	case TypeTimestamp,
		TypeDateTime,
		TypeDate:
		return scanTimestamp(src, dest, e.Loc)
	case TypeTinyInt, TypeSmallInt, TypeMediumInt, TypeInt, TypeBigInt, TypeTime, TypeYear:
		return scanInt64(src, dest)
	case TypeFloat, TypeDouble, TypeReal:
		return scanFloat(src, dest)
	case TypeChar, TypeVarChar, TypeTinyText, TypeText, TypeMediumText, TypeLongText:
		return scanString(src, dest)
	case TypeBoolean:
		return scanBool(src, dest)
	case TypeBinary, TypeVarBinary, TypeTinyBlob, TypeBlob, TypeMediumBlob, TypeLongBlob:
		// I suspect there might be some hex encoding
		return scanBinary(src, dest)
	case TypeEnum, TypeSet:
		return scanEnum(src, dest)
	case TypeGeometry, TypePoint, TypeLineString, TypePolygon, TypeMultiPoint, TypeMultiLineString,
		TypeMultiPolygon, TypeGeometryCollection:
		return scanGeometry(src, dest)
	case TypeDecimal, TypeNumeric:
		return scanDecimal(src, dest)
	}
	return fmt.Errorf("unsupported type %s", name)
}

func hasOptions(opts []TypeOption, opt TypeOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}
