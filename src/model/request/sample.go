package request

import "time"

type SampleModel struct {
	// サンプル文字列_必須
	SampleStringRequired string
	// サンプル文字列_最小文字数
	SampleStringMin string
	// サンプル文字列_最大文字数
	SampleStringMax string
	// サンプル文字列_範囲指定文字数
	SampleStringLength string
	// サンプル文字列_正規表現
	SampleStringRegexp string
	// サンプル文字列_正規表現_is使用
	SampleStringRegexpIs string
	// サンプル整数
	SampleInt int
	// サンプル整数8
	SampleInt8 int8
	// サンプル整数16
	SampleInt16 int16
	// サンプル整数32
	SampleInt32 int32
	// サンプル整数64
	SampleInt64 int64
	// サンプル符号無し整数
	SampleUint uint
	// サンプル符号無し整数8
	SampleUint8 uint8
	// サンプル符号無し整数16
	SampleUint16 uint16
	// サンプル符号無し整数32
	SampleUint32 uint32
	// サンプル符号無し整数64
	SampleUint64 uint64
	// サンプル小数32
	SampleFloat32 float32
	// サンプル小数64
	SampleFloat64 float64
	// サンプルバイト(エイリアス)
	SampleByte byte
	// サンプル真偽値
	SampleBool bool
	// サンプル日付
	SampleDate time.Time
	// サンプル配列
	SampleArray []string
	// サンプル構造体
	SampleSubStruct SampleSubModel
	// サンプル相関
	SampleStringRelation string
}

type SampleSubModel struct {
	// サンプル文字列_必須
	SampleStringRequired string
}
