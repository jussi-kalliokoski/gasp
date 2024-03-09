package token

import (
	"fmt"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected []Token
	}{
		{
			name: "basic literals",
			source: `
(literal 123 123.456 "foo" \q)
			`,
			expected: []Token{
				newToken(KindOpenParen, 1),
				newToken(KindSymbol, 7),
				newToken(KindWhitespace, 1),
				newInteger(3, BaseDecimal),
				newToken(KindWhitespace, 1),
				newFloat(7),
				newToken(KindWhitespace, 1),
				newString(5),
				newToken(KindWhitespace, 1),
				newCharacter(2),
				newToken(KindCloseParen, 1),
			},
		},
		{
			name: "binary integer",
			source: `
0b0101_1100_0111_0000
			`,
			expected: []Token{
				newInteger(21, BaseBinary),
			},
		},
		{
			name: "padded binary integer",
			source: `
0b_10_
			`,
			expected: []Token{
				newInteger(6, BaseBinary),
			},
		},
		{
			name: "empty binary integer",
			source: `
0b
			`,
			expected: []Token{
				newInteger(2, BaseBinary, literalFlagEmptyInt),
			},
		},
		{
			name: "invalid binary integer",
			source: `
0b012
			`,
			expected: []Token{
				newInteger(4, BaseBinary),
				newInteger(1, BaseDecimal),
			},
		},
		{
			name: "octal integer",
			source: `
0o7654_3210
			`,
			expected: []Token{
				newInteger(11, BaseOctal),
			},
		},
		{
			name: "padded octal integer",
			source: `
0o_01234567_
			`,
			expected: []Token{
				newInteger(12, BaseOctal),
			},
		},
		{
			name: "empty octal integer",
			source: `
0o
			`,
			expected: []Token{
				newInteger(2, BaseOctal, literalFlagEmptyInt),
			},
		},
		{
			name: "invalid octal integer",
			source: `
0o012345678
			`,
			expected: []Token{
				newInteger(10, BaseOctal),
				newInteger(1, BaseDecimal),
			},
		},
		{
			name: "decimal integer",
			source: `
1_234_567_890
			`,
			expected: []Token{
				newInteger(13, BaseDecimal),
			},
		},
		{
			name: "padded decimal integer",
			source: `
1234567890_
			`,
			expected: []Token{
				newInteger(11, BaseDecimal),
			},
		},
		{
			name: "zero-prefixed decimal integer",
			source: `
01234567890
			`,
			expected: []Token{
				newInteger(11, BaseDecimal),
			},
		},
		{
			name: "invalid decimal integer",
			source: `
1234567890A
			`,
			expected: []Token{
				newInteger(10, BaseDecimal),
				newToken(KindSymbol, 1),
			},
		},
		{
			name: "hexadecimal integer",
			source: `
0x0123456789AaBbCcDdEeFf
			`,
			expected: []Token{
				newInteger(24, BaseHexadecimal),
			},
		},
		{
			name: "padded hexadecimal integer",
			source: `
0x_0123456789AaBbCcDdEeFf_
			`,
			expected: []Token{
				newInteger(26, BaseHexadecimal),
			},
		},
		{
			name: "empty hexadecimal integer",
			source: `
0x
			`,
			expected: []Token{
				newInteger(2, BaseHexadecimal, literalFlagEmptyInt),
			},
		},
		{
			name: "invalid hexadecimal integer",
			source: `
0x0123456789AaBbCcDdEeFfGg
			`,
			expected: []Token{
				newInteger(24, BaseHexadecimal),
				newToken(KindSymbol, 2),
			},
		},
		{
			name: "valid floats",
			source: `
0123456789.0
12.3E10
98_76_._54e+592
00.444441e-123
0.9191e0
0.0000
1_234E56
			`,
			expected: []Token{
				newFloat(12),
				newToken(KindWhitespace, 1),
				newFloat(7),
				newToken(KindWhitespace, 1),
				newFloat(15),
				newToken(KindWhitespace, 1),
				newFloat(14),
				newToken(KindWhitespace, 1),
				newFloat(8),
				newToken(KindWhitespace, 1),
				newFloat(6),
				newToken(KindWhitespace, 1),
				newFloat(8),
			},
		},
		{
			name: "invalid binary float",
			source: `
0b1.0
			`,
			expected: []Token{
				newInteger(3, BaseBinary),
				newToken(KindSymbol, 2),
			},
		},
		{
			name: "invalid octal float",
			source: `
0o1.0
			`,
			expected: []Token{
				newInteger(3, BaseOctal),
				newToken(KindSymbol, 2),
			},
		},
		{
			name: "invalid hexadecimal float",
			source: `
0x1.0
			`,
			expected: []Token{
				newInteger(3, BaseHexadecimal),
				newToken(KindSymbol, 2),
			},
		},
		{
			name: "invalid dotted decimal",
			source: `
0.
			`,
			expected: []Token{
				newInteger(1, BaseDecimal),
				newToken(KindSymbol, 1),
			},
		},
		{
			name: "multiline string",
			source: `
"
hello
"
			`,
			expected: []Token{
				newString(9),
			},
		},
		{
			name: "string escapes",
			source: `
"\\\u3245\""
			`,
			expected: []Token{
				newString(12),
			},
		},
		{
			name: "unterminated string",
			source: `
"
			`,
			expected: []Token{
				newString(1, literalFlagUnterminated),
			},
		},
		{
			name: "valid characters",
			source: `
\\
\a
\1
\@
\(
\auml
\u1234
\:hello:
			`,
			expected: []Token{
				newCharacter(2),
				newToken(KindWhitespace, 1),
				newCharacter(2),
				newToken(KindWhitespace, 1),
				newCharacter(2),
				newToken(KindWhitespace, 1),
				newCharacter(2),
				newToken(KindWhitespace, 1),
				newCharacter(2),
				newToken(KindWhitespace, 1),
				newCharacter(5),
				newToken(KindWhitespace, 1),
				newCharacter(6),
				newToken(KindWhitespace, 1),
				newCharacter(8),
			},
		},
		{
			name:   "missing character",
			source: `\`,
			expected: []Token{
				newCharacter(1, literalFlagMissingChar),
			},
		},
		{
			name: "non-symbol after character",
			source: `
\@@
			`,
			expected: []Token{
				newCharacter(2),
				newToken(KindDeref, 1),
			},
		},
		{
			name: "valid symbols",
			source: `
:l33T_häxÖrZ:
.vAv.
a/b.c
+-*/!_?<>=:.
a'
			`,
			expected: []Token{
				newToken(KindSymbol, 15),
				newToken(KindWhitespace, 1),
				newToken(KindSymbol, 5),
				newToken(KindWhitespace, 1),
				newToken(KindSymbol, 5),
				newToken(KindWhitespace, 1),
				newToken(KindSymbol, 12),
				newToken(KindWhitespace, 1),
				newToken(KindSymbol, 2),
			},
		},
		{
			name: "line comments",
			source: `
;comment wee
(this ; this
	(is; is
	    (lisp); (comment)
	); hsssss
) ;
; ; ; ;
			`,
			expected: []Token{
				newToken(KindLineComment, 12),
				newToken(KindWhitespace, 1),
				newToken(KindOpenParen, 1),
				newToken(KindSymbol, 4),
				newToken(KindWhitespace, 1),
				newToken(KindLineComment, 6),
				newToken(KindWhitespace, 2),
				newToken(KindOpenParen, 1),
				newToken(KindSymbol, 2),
				newToken(KindLineComment, 4),
				newToken(KindWhitespace, 6),
				newToken(KindOpenParen, 1),
				newToken(KindSymbol, 4),
				newToken(KindCloseParen, 1),
				newToken(KindLineComment, 11),
				newToken(KindWhitespace, 2),
				newToken(KindCloseParen, 1),
				newToken(KindLineComment, 8),
				newToken(KindWhitespace, 1),
				newToken(KindCloseParen, 1),
				newToken(KindWhitespace, 1),
				newToken(KindLineComment, 1),
				newToken(KindWhitespace, 1),
				newToken(KindLineComment, 7),
			},
		},
		{
			name: "number properties",
			source: `
0.o
0.0.o
			`,
			expected: []Token{
				newInteger(1, BaseDecimal),
				newToken(KindSymbol, 2),
				newToken(KindWhitespace, 1),
				newFloat(3),
				newToken(KindSymbol, 2),
			},
		},
		{
			name:   "unquote",
			source: `~`,
			expected: []Token{
				newToken(KindUnquote, 1),
			},
		},
		{
			name:   "unquote splicing",
			source: `~@`,
			expected: []Token{
				newToken(KindUnquoteSplicing, 2),
			},
		},
		{
			name:   "quote",
			source: `'`,
			expected: []Token{
				newToken(KindQuote, 1),
			},
		},
		{
			name:   "backquote",
			source: "`",
			expected: []Token{
				newToken(KindBackquote, 1),
			},
		},
		{
			name:   "deref",
			source: "@",
			expected: []Token{
				newToken(KindDeref, 1),
			},
		},
		{
			name:   "metadata",
			source: "^",
			expected: []Token{
				newToken(KindMetadata, 1),
			},
		},
		{
			name:   "dispatch",
			source: "#",
			expected: []Token{
				newToken(KindDispatch, 1),
			},
		},
		{
			name: "groups",
			source: `
({[]})
			`,
			expected: []Token{
				newToken(KindOpenParen, 1),
				newToken(KindOpenBrace, 1),
				newToken(KindOpenBracket, 1),
				newToken(KindCloseBracket, 1),
				newToken(KindCloseBrace, 1),
				newToken(KindCloseParen, 1),
			},
		},
		{
			name:   "invalid token",
			source: "\u0000",
			expected: []Token{
				newToken(KindInvalid, 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc sliceConsumer
			source := strings.TrimSpace(tt.source)

			err := Tokenize(&sc, source)
			received := sc.Tokens()

			if err != nil {
				t.Fatal(err)
			}

			diffTokens(t, source, tt.expected, received)
		})
	}
}

func TestFormat(t *testing.T) {
	requireEqual(t, "token.Token{Kind:invalid Len:1}", fmt.Sprintf("%+v", newToken(KindInvalid, 1)))
	requireEqual(t, "{invalid 1}", fmt.Sprintf("%v", newToken(KindInvalid, 1)))
	requireEqual(t, "{invalid 1}", fmt.Sprintf("%s", newToken(KindInvalid, 1)))
	requireEqual(t, "token.Token{Kind:open-paren Len:1}", fmt.Sprintf("%+v", newToken(KindOpenParen, 1)))
	requireEqual(t, "{open-paren 1}", fmt.Sprintf("%v", newToken(KindOpenParen, 1)))
	requireEqual(t, "{open-paren 1}", fmt.Sprintf("%s", newToken(KindOpenParen, 1)))
	requireEqual(t, "token.Token{Kind:close-paren Len:1}", fmt.Sprintf("%+v", newToken(KindCloseParen, 1)))
	requireEqual(t, "{close-paren 1}", fmt.Sprintf("%v", newToken(KindCloseParen, 1)))
	requireEqual(t, "{close-paren 1}", fmt.Sprintf("%s", newToken(KindCloseParen, 1)))
	requireEqual(t, "token.Token{Kind:open-brace Len:1}", fmt.Sprintf("%+v", newToken(KindOpenBrace, 1)))
	requireEqual(t, "{open-brace 1}", fmt.Sprintf("%v", newToken(KindOpenBrace, 1)))
	requireEqual(t, "{open-brace 1}", fmt.Sprintf("%s", newToken(KindOpenBrace, 1)))
	requireEqual(t, "token.Token{Kind:close-brace Len:1}", fmt.Sprintf("%+v", newToken(KindCloseBrace, 1)))
	requireEqual(t, "{close-brace 1}", fmt.Sprintf("%v", newToken(KindCloseBrace, 1)))
	requireEqual(t, "{close-brace 1}", fmt.Sprintf("%s", newToken(KindCloseBrace, 1)))
	requireEqual(t, "token.Token{Kind:open-bracket Len:1}", fmt.Sprintf("%+v", newToken(KindOpenBracket, 1)))
	requireEqual(t, "{open-bracket 1}", fmt.Sprintf("%v", newToken(KindOpenBracket, 1)))
	requireEqual(t, "{open-bracket 1}", fmt.Sprintf("%s", newToken(KindOpenBracket, 1)))
	requireEqual(t, "token.Token{Kind:close-bracket Len:1}", fmt.Sprintf("%+v", newToken(KindCloseBracket, 1)))
	requireEqual(t, "{close-bracket 1}", fmt.Sprintf("%v", newToken(KindCloseBracket, 1)))
	requireEqual(t, "{close-bracket 1}", fmt.Sprintf("%s", newToken(KindCloseBracket, 1)))
	requireEqual(t, "token.Token{Kind:quote Len:1}", fmt.Sprintf("%+v", newToken(KindQuote, 1)))
	requireEqual(t, "{quote 1}", fmt.Sprintf("%v", newToken(KindQuote, 1)))
	requireEqual(t, "{quote 1}", fmt.Sprintf("%s", newToken(KindQuote, 1)))
	requireEqual(t, "token.Token{Kind:backquote Len:1}", fmt.Sprintf("%+v", newToken(KindBackquote, 1)))
	requireEqual(t, "{backquote 1}", fmt.Sprintf("%v", newToken(KindBackquote, 1)))
	requireEqual(t, "{backquote 1}", fmt.Sprintf("%s", newToken(KindBackquote, 1)))
	requireEqual(t, "token.Token{Kind:deref Len:1}", fmt.Sprintf("%+v", newToken(KindDeref, 1)))
	requireEqual(t, "{deref 1}", fmt.Sprintf("%v", newToken(KindDeref, 1)))
	requireEqual(t, "{deref 1}", fmt.Sprintf("%s", newToken(KindDeref, 1)))
	requireEqual(t, "token.Token{Kind:metadata Len:1}", fmt.Sprintf("%+v", newToken(KindMetadata, 1)))
	requireEqual(t, "{metadata 1}", fmt.Sprintf("%v", newToken(KindMetadata, 1)))
	requireEqual(t, "{metadata 1}", fmt.Sprintf("%s", newToken(KindMetadata, 1)))
	requireEqual(t, "token.Token{Kind:dispatch Len:1}", fmt.Sprintf("%+v", newToken(KindDispatch, 1)))
	requireEqual(t, "{dispatch 1}", fmt.Sprintf("%v", newToken(KindDispatch, 1)))
	requireEqual(t, "{dispatch 1}", fmt.Sprintf("%s", newToken(KindDispatch, 1)))
	requireEqual(t, "token.Token{Kind:unquote Len:1}", fmt.Sprintf("%+v", newToken(KindUnquote, 1)))
	requireEqual(t, "{unquote 1}", fmt.Sprintf("%v", newToken(KindUnquote, 1)))
	requireEqual(t, "{unquote 1}", fmt.Sprintf("%s", newToken(KindUnquote, 1)))
	requireEqual(t, "token.Token{Kind:unquote-splicing Len:2}", fmt.Sprintf("%+v", newToken(KindUnquoteSplicing, 2)))
	requireEqual(t, "{unquote-splicing 2}", fmt.Sprintf("%v", newToken(KindUnquoteSplicing, 2)))
	requireEqual(t, "{unquote-splicing 2}", fmt.Sprintf("%s", newToken(KindUnquoteSplicing, 2)))
	requireEqual(t, "token.Token{Kind:symbol Len:3}", fmt.Sprintf("%+v", newToken(KindSymbol, 3)))
	requireEqual(t, "{symbol 3}", fmt.Sprintf("%v", newToken(KindSymbol, 3)))
	requireEqual(t, "{symbol 3}", fmt.Sprintf("%s", newToken(KindSymbol, 3)))
	requireEqual(t, "token.Token{Kind:whitespace Len:5}", fmt.Sprintf("%+v", newToken(KindWhitespace, 5)))
	requireEqual(t, "{whitespace 5}", fmt.Sprintf("%v", newToken(KindWhitespace, 5)))
	requireEqual(t, "{whitespace 5}", fmt.Sprintf("%s", newToken(KindWhitespace, 5)))
	requireEqual(t, "token.Token{Kind:line-comment Len:4}", fmt.Sprintf("%+v", newToken(KindLineComment, 4)))
	requireEqual(t, "{line-comment 4}", fmt.Sprintf("%v", newToken(KindLineComment, 4)))
	requireEqual(t, "{line-comment 4}", fmt.Sprintf("%s", newToken(KindLineComment, 4)))
	requireEqual(t, "token.Token{Kind:token.Literal{token.Integer{Base:10 EmptyInt:false}} Len:12}", fmt.Sprintf("%+v", newInteger(12, BaseDecimal)))
	requireEqual(t, "token.Token{Kind:token.Literal{token.Integer{Base:16 EmptyInt:true}} Len:12}", fmt.Sprintf("%+v", newInteger(12, BaseHexadecimal, literalFlagEmptyInt)))
	requireEqual(t, "{{token.Integer{10 false}} 12}", fmt.Sprintf("%v", newInteger(12, BaseDecimal)))
	requireEqual(t, "{token.Literal 12}", fmt.Sprintf("%s", newInteger(12, BaseDecimal)))
	requireEqual(t, "{token.Integer}", fmt.Sprintf("%s", newInteger(12, BaseDecimal).Literal()))
	requireEqual(t, "token.Token{Kind:token.Literal{token.Float{EmptyExponent:false}} Len:11}", fmt.Sprintf("%+v", newFloat(11)))
	requireEqual(t, "token.Token{Kind:token.Literal{token.Float{EmptyExponent:true}} Len:11}", fmt.Sprintf("%+v", newFloat(11, literalFlagEmptyExponent)))
	requireEqual(t, "{{token.Float{false}} 11}", fmt.Sprintf("%v", newFloat(11)))
	requireEqual(t, "{token.Literal 11}", fmt.Sprintf("%s", newFloat(11)))
	requireEqual(t, "{token.Float}", fmt.Sprintf("%s", newFloat(11).Literal()))
	requireEqual(t, "token.Token{Kind:token.Literal{token.String{Unterminated:false}} Len:11}", fmt.Sprintf("%+v", newString(11)))
	requireEqual(t, "token.Token{Kind:token.Literal{token.String{Unterminated:true}} Len:11}", fmt.Sprintf("%+v", newString(11, literalFlagUnterminated)))
	requireEqual(t, "{{token.String{false}} 11}", fmt.Sprintf("%v", newString(11)))
	requireEqual(t, "{token.Literal 11}", fmt.Sprintf("%s", newString(11)))
	requireEqual(t, "{token.String}", fmt.Sprintf("%s", newString(11).Literal()))
	requireEqual(t, "token.Token{Kind:token.Literal{token.Character{MissingCharacter:false}} Len:2}", fmt.Sprintf("%+v", newCharacter(2)))
	requireEqual(t, "token.Token{Kind:token.Literal{token.Character{MissingCharacter:true}} Len:1}", fmt.Sprintf("%+v", newCharacter(1, literalFlagUnterminated)))
	requireEqual(t, "{{token.Character{false}} 2}", fmt.Sprintf("%v", newCharacter(2)))
	requireEqual(t, "{token.Literal 2}", fmt.Sprintf("%s", newCharacter(2)))
	requireEqual(t, "{token.Character}", fmt.Sprintf("%s", newCharacter(2).Literal()))
}

func TestPanics(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{
			"unrecognized token kind",
			func() {
				newToken(^Kind(0), 1).Format(&formatState{}, 'v')
			},
		},
		{
			"Literal called on non-literal",
			func() {
				_ = newToken(KindOpenParen, 1).Literal()
			},
		},
		{
			"String called on non-string literal",
			func() {
				_ = newInteger(1, BaseDecimal).Literal().String()
			},
		},
		{
			"Integer called on non-integer literal",
			func() {
				_ = newFloat(2).Literal().Integer()
			},
		},
		{
			"Float called on non-float literal",
			func() {
				_ = newCharacter(2).Literal().Float()
			},
		},
		{
			"Character called on non-character literal",
			func() {
				_ = newString(2).Literal().Character()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if p := recover(); p == nil {
					t.Fatal("expected a panic")
				}
			}()
			tt.fn()
		})
	}
}

type sliceConsumer struct {
	tokens []Token
}

func (sc *sliceConsumer) ConsumeToken(t Token) {
	sc.tokens = append(sc.tokens, t)
}

func (sc *sliceConsumer) Tokens() []Token {
	return sc.tokens
}

func newToken(kind Kind, tlen int) Token {
	return Token{
		kind: kind,
		len:  uint32(tlen),
	}
}

func newInteger(tlen int, base Base, flags ...literalFlags) Token {
	return Token{
		kind: KindLiteral,
		len:  uint32(tlen),
		literal: Literal{
			kind:  LiteralKindInteger,
			base:  base,
			flags: combineFlags(flags),
		},
	}
}

func newFloat(tlen int, flags ...literalFlags) Token {
	return Token{
		kind: KindLiteral,
		len:  uint32(tlen),
		literal: Literal{
			kind:  LiteralKindFloat,
			flags: combineFlags(flags),
		},
	}
}

func newString(tlen int, flags ...literalFlags) Token {
	return Token{
		kind: KindLiteral,
		len:  uint32(tlen),
		literal: Literal{
			kind:  LiteralKindString,
			flags: combineFlags(flags),
		},
	}
}

func newCharacter(tlen int, flags ...literalFlags) Token {
	return Token{
		kind: KindLiteral,
		len:  uint32(tlen),
		literal: Literal{
			kind:  LiteralKindCharacter,
			flags: combineFlags(flags),
		},
	}
}

func combineFlags(flags []literalFlags) literalFlags {
	var result literalFlags
	for _, f := range flags {
		result |= f
	}
	return result
}

func diffTokens(tb testing.TB, source string, a, b []Token) {
	pos := 0
	for i := 0; i < len(a); i++ {
		expected := a[i]
		if i >= len(b) {
			debugLine(tb, source, pos)
			tb.Fatalf("expected %+v, received EOF", expected)
		}

		received := b[i]
		if expected != received {
			debugLine(tb, source, pos)
			tb.Fatalf("expected %+v, received %+v", expected, received)
		}
		pos += expected.Len()
	}

	if len(b) > len(a) {
		received := b[len(a)]
		debugLine(tb, source, pos)
		tb.Fatalf("expected EOF, received %+v", received)
	}
}

func debugLine(tb testing.TB, source string, pos int) {
	line := 0
	col := pos
	for {
		lineEnd := strings.IndexRune(source, '\n')
		if lineEnd != -1 && col >= lineEnd {
			source = source[lineEnd+1:]
			col -= lineEnd + 1
			line++
			continue
		}

		currentLine := source
		if lineEnd != -1 {
			currentLine = currentLine[:lineEnd]
		}
		prefix := fmt.Sprintf("%d:%d: ", line+1, col+1)
		tb.Logf("%s%s", prefix, currentLine)
		tb.Logf("%s^", strings.Repeat(" ", col+len(prefix)))
		break
	}
}

func requireEqual[T comparable](tb testing.TB, expected, received T) {
	tb.Helper()
	if expected != received {
		tb.Fatalf("expected %v, received %v", expected, received)
	}
}

type formatState struct {
	flags   []int
	b       []byte
	wid     int
	widSet  bool
	prec    int
	precSet bool
}

func (s *formatState) Write(b []byte) (n int, err error) {
	s.b = append(s.b, b...)
	n = len(b)
	return
}

func (s *formatState) Width() (wid int, ok bool) {
	return s.wid, s.widSet
}

func (s *formatState) Precision() (prec int, ok bool) {
	return s.prec, s.precSet
}

func (s *formatState) Flag(c int) bool {
	for _, f := range s.flags {
		if f == c {
			return true
		}
	}
	return false
}
