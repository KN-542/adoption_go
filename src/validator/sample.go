package validator

import (
	"api/src/model"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ISampleValidator interface {
	Sample(m *model.SampleModel) error
	SampleSub(m *model.SampleSubModel) error
}

type SampleValidator struct{}

func NewSampleValidator() ISampleValidator {
	return &SampleValidator{}
}

func (v *SampleValidator) Sample(m *model.SampleModel) error {
	return validation.ValidateStruct(
		m,
		// サンプル文字列_必須
		validation.Field(
			&m.SampleStringRequired,
			validation.Required,
		),
		// サンプル文字列_最小文字数
		validation.Field(
			&m.SampleStringMin,
			validation.Length(3, 0),
		),
		// サンプル文字列_最大文字数
		validation.Field(
			&m.SampleStringMax,
			validation.Length(0, 7),
		),
		// サンプル文字列_範囲指定文字数
		validation.Field(
			&m.SampleStringLength,
			validation.Length(3, 7),
		),
		// サンプル文字列_正規表現
		validation.Field(
			&m.SampleStringRegexp,
			validation.Match(regexp.MustCompile(`^[0-9a-zA-Z]*$`)),
		),
		// サンプル文字列_正規表現_is使用
		validation.Field(
			&m.SampleStringRegexpIs,
			is.Alphanumeric,
		),
		// サンプル整数
		validation.Field(
			&m.SampleInt,
			validation.Min(int(-5)),
			validation.Max(int(5)),
		),
		// サンプル整数8
		validation.Field(
			&m.SampleInt8,
			validation.Min(int8(-5)),
			validation.Max(int8(5)),
		),
		// サンプル整数16
		validation.Field(
			&m.SampleInt16,
			validation.Min(int16(-5)),
			validation.Max(int16(5)),
		),
		// サンプル整数32
		validation.Field(
			&m.SampleInt32,
			validation.Min(int32(-5)),
			validation.Max(int32(5)),
		),
		// サンプル整数64
		validation.Field(
			&m.SampleInt64,
			validation.Min(int64(-5)),
			validation.Max(int64(5)),
		),
		// サンプル符号無し整数
		validation.Field(
			&m.SampleUint,
			validation.Min(uint(5)),
			validation.Max(uint(10)),
		),
		// サンプル符号無し整数8
		validation.Field(
			&m.SampleUint8,
			validation.Min(uint(5)),
			validation.Max(uint(10)),
		),
		// サンプル符号無し整数16
		validation.Field(
			&m.SampleUint16,
			validation.Min(uint(5)),
			validation.Max(uint(10)),
		),
		// サンプル符号無し整数32
		validation.Field(
			&m.SampleUint32,
			validation.Min(uint(5)),
			validation.Max(uint(10)),
		),
		// サンプル符号無し整数64
		validation.Field(
			&m.SampleUint64,
			validation.Min(uint(5)),
			validation.Max(uint(10)),
		),
		// サンプル配列
		validation.Field(
			&m.SampleArray,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(validation.Required),
		),
		// サンプル相関
		validation.Field(
			&m.SampleStringRelation,
			validation.By(func(value interface{}) error {
				return IsEqualStrings(&m.SampleStringRelation, &m.SampleStringRequired)
			}),
		),
	)
}

func (v *SampleValidator) SampleSub(m *model.SampleSubModel) error {
	return validation.ValidateStruct(
		m,
		// サンプル文字列_必須
		validation.Field(
			&m.SampleStringRequired,
			validation.Required,
		),
	)
}
