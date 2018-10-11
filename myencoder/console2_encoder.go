package myencoder

import (
	"go.uber.org/zap/zapcore"
	"time"
	"go.uber.org/zap/buffer"
	"unicode/utf8"
	"sync"
	"math"
	"golab/bufferpool"
)

// For JSON-escaping; see jsonEncoder.safeAddString below.
const _hex = "0123456789abcdef"

var _myPool = sync.Pool{New: func() interface{} {
	return &Console2Encoder{}
}}

func getConsole2Encoder() *Console2Encoder {
	return _myPool.Get().(*Console2Encoder)
}


func putConsole2Encoder(enc *Console2Encoder) {
	//if enc.reflectBuf != nil {
	//	enc.reflectBuf.Free()
	//}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.spaced = false
	enc.openNamespaces = 0
	//enc.reflectBuf = nil
	//enc.reflectEnc = nil
	_myPool.Put(enc)
}

// **注意** 需要实现zapcore.Encoder 所定义的接口
// **注意** 需要实现zapcore.PrimitiveArrayEncoder 所定义的接口
type Console2Encoder struct {
	*zapcore.EncoderConfig
	buf    *buffer.Buffer
	spaced bool // include spaces after colons and commas
	openNamespaces int
}


func NewConsole2Encoder(cfg zapcore.EncoderConfig, spaced bool) *Console2Encoder{
	return &Console2Encoder{
		EncoderConfig: &cfg,
		buf:           bufferpool.Get(),
		spaced:        spaced,
	}
}
//// AddReflected uses reflection to serialize arbitrary objects, so it's slow
//// and allocation-heavy.
//AddReflected(key string, value interface{}) error
//// OpenNamespace opens an isolated namespace where all subsequent fields will
//// be added. Applications can use namespaces to prevent key collisions when
//// injecting loggers into sub-components or third-party libraries.
//OpenNamespace(key string)

// Clone copies the encoder, ensuring that adding fields to the copy doesn't
//// affect the original.
//Clone() Encoder
//
//// EncodeEntry encodes an entry and fields, along with any accumulated
//// context, into a byte buffer and returns it.
//EncodeEntry(Entry, []Field) (*buffer.Buffer, error)

func (encoder Console2Encoder) AddReflected(key string, value interface{}) error {
	return nil
}

func (encoder Console2Encoder) OpenNamespace(key string) {}

func (encoder Console2Encoder) Clone() zapcore.Encoder {
	clone := encoder.clone()
	clone.buf.Write(encoder.buf.Bytes())
	return clone
}


func (encoder Console2Encoder) clone() *Console2Encoder{
	clone := getConsole2Encoder()
	clone.EncoderConfig = encoder.EncoderConfig
	clone.spaced = encoder.spaced
	clone.openNamespaces = encoder.openNamespaces
	clone.buf = bufferpool.Get()
	return clone

}




func (encoder Console2Encoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	return nil
}

func (encoder Console2Encoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	return nil
}

func (encoder Console2Encoder) AddBinary(key string, value []byte) {

}

func (encoder Console2Encoder) AddByteString(key string, value []byte) {
	encoder.addKey(key)
	encoder.AppendByteString(value)
}

func (enc *Console2Encoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}


// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *Console2Encoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

func (encoder Console2Encoder) AddBool(key string, value bool) {
	encoder.addKey(key)
	encoder.AppendBool(value)
}

func (enc *Console2Encoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (encoder Console2Encoder) AddComplex128(key string, value complex128) {
	encoder.addKey(key)
	encoder.AppendComplex128(value)
}

func (enc *Console2Encoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}


func (encoder Console2Encoder) AddComplex64(key string, value complex64) {

}
func (encoder Console2Encoder) AddDuration(key string, value time.Duration) {

}

func (encoder Console2Encoder) AddFloat64(key string, value float64) {

}

func (encoder Console2Encoder) AddFloat32(key string, value float32) {

}

func (encoder Console2Encoder) AddInt(key string, value int) {

}

func (encoder Console2Encoder) AddInt64(key string, value int64) {
	encoder.addKey(key)
	encoder.AppendInt64(value)
}

func (enc *Console2Encoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *Console2Encoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
		if enc.spaced {
			enc.buf.AppendByte(' ')
		}
	}
}

func (enc *Console2Encoder) addKey(key string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(key)
	enc.buf.AppendByte('"')
	enc.buf.AppendByte(':')
	if enc.spaced {
		enc.buf.AppendByte(' ')
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *Console2Encoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *Console2Encoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *Console2Encoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func (encoder Console2Encoder) AddInt32(key string, value int32) {

}

func (encoder Console2Encoder) AddInt16(key string, value int16) {

}
func (encoder Console2Encoder) AddInt8(key string, value int8) {

}

func (encoder Console2Encoder) AddString(key string, value string) {
	encoder.addKey(key)
	encoder.AppendString(value)
}

func (enc *Console2Encoder) AppendString(val string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(val)
	enc.buf.AppendByte('"')
}

func (encoder Console2Encoder) AddTime(key string, value time.Time) {

}

func (encoder Console2Encoder) AddUint(key string, value uint) {

}

func (encoder Console2Encoder) AddUint64(key string, value uint64) {

}

func (encoder Console2Encoder) AddUint32(key string, value uint32) {

}

func (encoder Console2Encoder) AddUint16(key string, value uint16) {

}

func (encoder Console2Encoder) AddUint8(key string, value uint8) {

}

func (encoder Console2Encoder) AddUintptr(key string, value uintptr) {

}

func (enc *Console2Encoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *Console2Encoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *Console2Encoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *Console2Encoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *Console2Encoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *Console2Encoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *Console2Encoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *Console2Encoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *Console2Encoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *Console2Encoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *Console2Encoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *Console2Encoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *Console2Encoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *Console2Encoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

func addFields(enc zapcore.ObjectEncoder, fields [] zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}


func (enc *Console2Encoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
}


func (enc Console2Encoder) EncodeEntry(ent zapcore.Entry, fields [] zapcore.Field) (*buffer.Buffer, error){
	final := enc.clone()
	// 去掉前面的大括号
	//final.buf.AppendByte('{')

	if final.LevelKey != "" {
		final.addKey(final.LevelKey)
		cur := final.buf.Len()
		final.EncodeLevel(ent.Level, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeLevel was a no-op. Fall back to strings to keep
			// output JSON valid.
			final.AppendString(ent.Level.String())
		}
	}
	if final.TimeKey != "" {
		final.AddTime(final.TimeKey, ent.Time)
	}
	if ent.LoggerName != "" && final.NameKey != "" {
		final.addKey(final.NameKey)
		cur := final.buf.Len()
		nameEncoder := final.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output JSON valid.
			final.AppendString(ent.LoggerName)
		}
	}
	if ent.Caller.Defined && final.CallerKey != "" {
		final.addKey(final.CallerKey)
		cur := final.buf.Len()
		final.EncodeCaller(ent.Caller, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeCaller was a no-op. Fall back to strings to
			// keep output JSON valid.
			final.AppendString(ent.Caller.String())
		}
	}
	if final.MessageKey != "" {
		final.addKey(enc.MessageKey)
		final.AppendString(ent.Message)
	}
	if enc.buf.Len() > 0 {
		//final.addElementSeparator()
		final.buf.Write(enc.buf.Bytes())
	}

	// 用$$$ 隔开field
	final.buf.AppendString("$$$ ");


	addFields(final, fields)
	final.closeOpenNamespaces()
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
	}

	// 去掉后面的大括号
	//final.buf.AppendByte('}')
	if final.LineEnding != "" {
		final.buf.AppendString(final.LineEnding)
	} else {
		final.buf.AppendString(zapcore.DefaultLineEnding)
	}

	ret := final.buf
	putConsole2Encoder(final)
	return ret, nil
}

