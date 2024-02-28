package validator

import (
	"api/src/model"
	"testing"
)

func TestSampleValidator_Sample(t *testing.T) {
	type args struct {
		m *model.SampleModel
	}
	tests := []struct {
		name    string
		v       *SampleValidator
		args    args
		wantErr bool
	}{
		// 必須 ok
		{
			"ok_required",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
			}},
			false,
		},
		// 必須 ng nil
		{
			"ng_required_nil",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleArray: []string{"a"},
			}},
			true,
		},
		// 必須 ng 空文字
		{
			"ng_required_empty",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "",
				SampleArray:          []string{"a"},
			}},
			true,
		},
		// 最小文字数 ok (最小文字数)
		{
			"ok_min_length",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringMin:      "abc",
			}},
			false,
		},
		// 最小文字数 ng (最小文字数 - 1)
		{
			"ng_min_length",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringMin:      "ab",
			}},
			true,
		},
		// 最大文字数 ok (最大文字数)
		{
			"ok_max_length",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringMax:      "abcdefg",
			}},
			false,
		},
		// 最大文字数 ng (最大文字数 + 1)
		{
			"ng_max_length",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringMax:      "abcdefgh",
			}},
			true,
		},
		// 範囲指定文字数 ok (最小文字数)
		{
			"ok_length_1",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringLength:   "abc",
			}},
			false,
		},
		// 範囲指定文字数 ok (最大文字数)
		{
			"ok_length_2",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringLength:   "abcdefg",
			}},
			false,
		},
		// 範囲指定文字数 ng (最小文字数 - 1)
		{
			"ng_length_1",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringLength:   "ab",
			}},
			true,
		},
		// 範囲指定文字数 ng (最大文字数 + 1)
		{
			"ng_length_2",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringLength:   "abcdefgh",
			}},
			true,
		},
		// 正規表現 ok
		{
			"ok_regexp",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringRegexp:   "a",
			}},
			false,
		},
		// 正規表現 ng
		{
			"ng_regexp",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringRegexp:   "@",
			}},
			true,
		},
		// 正規表現_is使用 ok
		{
			"ok_regexp_is",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringRegexpIs: "a",
			}},
			false,
		},
		// 正規表現_is使用 ng
		{
			"ng_regexp_is",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringRegexpIs: "@",
			}},
			true,
		},
		// サンプル整数_最小値 ok min
		{
			"int_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt:            -5,
			}},
			false,
		},
		// サンプル整数_最大値 ok max
		{
			"int_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt:            5,
			}},
			false,
		},
		// サンプル整数_最小値 ng min - 1
		{
			"int_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt:            -5 - 1,
			}},
			true,
		},
		// サンプル整数_最大値 ok max + 1
		{
			"int_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt:            5 + 1,
			}},
			true,
		},
		// サンプル整数8_最小値 ok min
		{
			"int8_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt8:           -5,
			}},
			false,
		},
		// サンプル整数8_最大値 ok max
		{
			"int8_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt8:           5,
			}},
			false,
		},
		// サンプル整数8_最小値 ng min - 1
		{
			"int8_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt8:           -5 - 1,
			}},
			true,
		},
		// サンプル整数8_最大値 ok max + 1
		{
			"int8_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt8:           5 + 1,
			}},
			true,
		},
		// サンプル整数16_最小値 ok min
		{
			"int16_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt16:          -5,
			}},
			false,
		},
		// サンプル整数16_最大値 ok max
		{
			"int16_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt16:          5,
			}},
			false,
		},
		// サンプル整数16_最小値 ng min - 1
		{
			"int16_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt16:          -5 - 1,
			}},
			true,
		},
		// サンプル整数16_最大値 ok max + 1
		{
			"int16_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt16:          5 + 1,
			}},
			true,
		},
		// サンプル整数32_最小値 ok min
		{
			"int32_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt32:          -5,
			}},
			false,
		},
		// サンプル整数32_最大値 ok max
		{
			"int32_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt32:          5,
			}},
			false,
		},
		// サンプル整数32_最小値 ng min - 1
		{
			"int32_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt32:          -5 - 1,
			}},
			true,
		},
		// サンプル整数32_最大値 ok max + 1
		{
			"int32_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt32:          5 + 1,
			}},
			true,
		},
		// サンプル整数64_最小値 ok min
		{
			"int64_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt64:          -5,
			}},
			false,
		},
		// サンプル整数64_最大値 ok max
		{
			"int64_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt64:          5,
			}},
			false,
		},
		// サンプル整数64_最小値 ng min - 1
		{
			"int64_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt64:          -5 - 1,
			}},
			true,
		},
		// サンプル整数64_最大値 ok max + 1
		{
			"int64_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleInt64:          5 + 1,
			}},
			true,
		},
		// サンプル符号無し整数 ok min
		{
			"uint_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint:           5,
			}},
			false,
		},
		// サンプル符号無し整数 ok max
		{
			"uint_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint:           10,
			}},
			false,
		},
		// サンプル符号無し整数 ng min - 1
		{
			"uint_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint:           5 - 1,
			}},
			true,
		},
		// サンプル符号無し整数 ng max + 1
		{
			"uint_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint:           10 + 1,
			}},
			true,
		},
		// サンプル符号無し整数8 ok min
		{
			"uint8_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint8:          5,
			}},
			false,
		},
		// サンプル符号無し整数8 ok max
		{
			"uint8_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint8:          10,
			}},
			false,
		},
		// サンプル符号無し整数8 ng min - 1
		{
			"uint8_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint8:          5 - 1,
			}},
			true,
		},
		// サンプル符号無し整数8 ng max + 1
		{
			"uint8_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint8:          10 + 1,
			}},
			true,
		},
		// サンプル符号無し整数16 ok min
		{
			"uint16_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint16:         5,
			}},
			false,
		},
		// サンプル符号無し整数16 ok max
		{
			"uint16_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint16:         10,
			}},
			false,
		},
		// サンプル符号無し整数16 ng min - 1
		{
			"uint16_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint16:         5 - 1,
			}},
			true,
		},
		// サンプル符号無し整数16 ng max + 1
		{
			"uint16_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint16:         10 + 1,
			}},
			true,
		},
		// サンプル符号無し整数32 ok min
		{
			"uint32_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint32:         5,
			}},
			false,
		},
		// サンプル符号無し整数32 ok max
		{
			"uint32_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint32:         10,
			}},
			false,
		},
		// サンプル符号無し整数32 ng min - 1
		{
			"uint32_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint32:         5 - 1,
			}},
			true,
		},
		// サンプル符号無し整数32 ng max + 1
		{
			"uint32_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint32:         10 + 1,
			}},
			true,
		},
		// サンプル符号無し整数64 ok min
		{
			"uint64_min_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint64:         5,
			}},
			false,
		},
		// サンプル符号無し整数64 ok max
		{
			"uint64_max_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint64:         10,
			}},
			false,
		},
		// サンプル符号無し整数64 ng min - 1
		{
			"uint64_min_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint64:         5 - 1,
			}},
			true,
		},
		// サンプル符号無し整数64 ng max + 1
		{
			"uint64_max_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleUint64:         10 + 1,
			}},
			true,
		},
		// サンプル配列 ng 必須
		{
			"array_ng_nil",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
			}},
			true,
		},
		// サンプル配列 ng 空配列
		{
			"array_ng_empty",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{},
			}},
			true,
		},
		// サンプル配列 ng 要素が空
		{
			"array_ok_el_empty",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{""},
			}},
			true,
		},
		// サンプル相関 ok
		{
			"relation_ok",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringRelation: "a",
			}},
			false,
		},
		// サンプル相関 ng
		{
			"relation_ng",
			&SampleValidator{},
			args{&model.SampleModel{
				SampleStringRequired: "a",
				SampleArray:          []string{"a"},
				SampleStringRelation: "ab",
			}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Sample(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SampleValidator.Sample() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSampleValidator_SampleSub(t *testing.T) {
	type args struct {
		m *model.SampleSubModel
	}
	tests := []struct {
		name    string
		v       *SampleValidator
		args    args
		wantErr bool
	}{
		// 必須 ok
		{
			"ok_required",
			&SampleValidator{},
			args{&model.SampleSubModel{
				SampleStringRequired: "a",
			}},
			false,
		},
		// 必須 ng nil
		{
			"ng_required_nil",
			&SampleValidator{},
			args{&model.SampleSubModel{}},
			true,
		},
		// 必須 ng 空文字
		{
			"ng_required_empty",
			&SampleValidator{},
			args{&model.SampleSubModel{
				SampleStringRequired: "",
			}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.SampleSub(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SampleValidator.SampleSub() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
