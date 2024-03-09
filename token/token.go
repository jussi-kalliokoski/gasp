package token

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	kind    Kind
	len     uint32
	literal Literal
}

func (t Token) Kind() Kind {
	return t.kind
}

func (t Token) Len() int {
	return int(t.len)
}

func (t Token) Literal() Literal {
	if k1, k2 := KindLiteral, t.kind; k1 != k2 {
		panic(fmt.Errorf("expected %v, got %v", k1, k2))
	}
	return t.literal
}

func (t Token) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			if t.Kind() == KindLiteral {
				fmt.Fprintf(f, "%T{Kind:%+v Len:%d}", t, t.Literal(), t.Len())
			} else {
				fmt.Fprintf(f, "%T{Kind:%s Len:%d}", t, t.Kind().string(), t.Len())
			}
		} else {
			if t.Kind() == KindLiteral {
				fmt.Fprintf(f, "{%v %d}", t.Literal(), t.Len())
			} else {
				fmt.Fprintf(f, "{%s %d}", t.Kind().string(), t.Len())
			}
		}
	default:
		if t.Kind() == KindLiteral {
			fmt.Fprintf(f, "{%T %d}", t.Literal(), t.Len())
		} else {
			fmt.Fprintf(f, "{%s %d}", t.Kind().string(), t.Len())
		}
	}
}

type Kind uint8

const (
	KindInvalid Kind = iota
	kindNone
	KindLiteral
	KindOpenParen
	KindCloseParen
	KindOpenBrace
	KindCloseBrace
	KindOpenBracket
	KindCloseBracket
	KindQuote
	KindBackquote
	KindDeref
	KindMetadata
	KindDispatch
	KindUnquote
	KindUnquoteSplicing
	KindSymbol
	KindWhitespace
	KindLineComment
)

func (k Kind) string() string {
	switch k {
	case KindInvalid:
		return "invalid"
	case KindOpenParen:
		return "open-paren"
	case KindCloseParen:
		return "close-paren"
	case KindOpenBrace:
		return "open-brace"
	case KindCloseBrace:
		return "close-brace"
	case KindOpenBracket:
		return "open-bracket"
	case KindCloseBracket:
		return "close-bracket"
	case KindQuote:
		return "quote"
	case KindBackquote:
		return "backquote"
	case KindDeref:
		return "deref"
	case KindMetadata:
		return "metadata"
	case KindDispatch:
		return "dispatch"
	case KindUnquote:
		return "unquote"
	case KindUnquoteSplicing:
		return "unquote-splicing"
	case KindSymbol:
		return "symbol"
	case KindWhitespace:
		return "whitespace"
	case KindLineComment:
		return "line-comment"
	default:
		panic(fmt.Errorf("unknown kind: %v", k))
	}
}

type Literal struct {
	kind  LiteralKind
	base  Base
	flags literalFlags
}

func (l Literal) Kind() LiteralKind {
	return l.kind
}

func (l Literal) Integer() Integer {
	if k1, k2 := LiteralKindInteger, l.Kind(); k1 != k2 {
		panic(fmt.Errorf("expected %v, got %v", k1, k2))
	}
	return Integer{l}
}

func (l Literal) Float() Float {
	if k1, k2 := LiteralKindFloat, l.Kind(); k1 != k2 {
		panic(fmt.Errorf("expected %v, got %v", k1, k2))
	}
	return Float{l}
}

func (l Literal) String() String {
	if k1, k2 := LiteralKindString, l.Kind(); k1 != k2 {
		panic(fmt.Errorf("expected %v, got %v", k1, k2))
	}
	return String{l}
}

func (l Literal) Character() Character {
	if k1, k2 := LiteralKindCharacter, l.Kind(); k1 != k2 {
		panic(fmt.Errorf("expected %v, got %v", k1, k2))
	}
	return Character{l}
}

func (l Literal) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			switch l.kind {
			case LiteralKindInteger:
				fmt.Fprintf(f, "%T{%+v}", l, l.Integer())
			case LiteralKindFloat:
				fmt.Fprintf(f, "%T{%+v}", l, l.Float())
			case LiteralKindString:
				fmt.Fprintf(f, "%T{%+v}", l, l.String())
			case LiteralKindCharacter:
				fmt.Fprintf(f, "%T{%+v}", l, l.Character())
			}
		} else {
			switch l.kind {
			case LiteralKindInteger:
				fmt.Fprintf(f, "{%T%v}", l.Integer(), l.Integer())
			case LiteralKindFloat:
				fmt.Fprintf(f, "{%T%v}", l.Float(), l.Float())
			case LiteralKindString:
				fmt.Fprintf(f, "{%T%v}", l.String(), l.String())
			case LiteralKindCharacter:
				fmt.Fprintf(f, "{%T%v}", l.Character(), l.Character())
			}
		}
	default:
		switch l.kind {
		case LiteralKindInteger:
			fmt.Fprintf(f, "{%T}", l.Integer())
		case LiteralKindFloat:
			fmt.Fprintf(f, "{%T}", l.Float())
		case LiteralKindString:
			fmt.Fprintf(f, "{%T}", l.String())
		case LiteralKindCharacter:
			fmt.Fprintf(f, "{%T}", l.Character())
		}
	}
}

type Integer struct {
	l Literal
}

func (l Integer) Base() Base {
	return l.l.base
}

func (l Integer) EmptyInt() bool {
	return l.l.flags&literalFlagEmptyInt != 0
}

func (l Integer) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('+') {
		fmt.Fprintf(f, "%T{Base:%d EmptyInt:%t}", l, l.Base(), l.EmptyInt())
	} else {
		fmt.Fprintf(f, "{%d %t}", l.Base(), l.EmptyInt())
	}
}

type Float struct {
	l Literal
}

func (l Float) EmptyExponent() bool {
	return l.l.flags&literalFlagEmptyExponent != 0
}

func (l Float) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('+') {
		fmt.Fprintf(f, "%T{EmptyExponent:%t}", l, l.EmptyExponent())
	} else {
		fmt.Fprintf(f, "{%t}", l.EmptyExponent())
	}
}

type String struct {
	l Literal
}

func (l String) Unterminated() bool {
	return l.l.flags&literalFlagUnterminated != 0
}

func (l String) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('+') {
		fmt.Fprintf(f, "%T{Unterminated:%t}", l, l.Unterminated())
	} else {
		fmt.Fprintf(f, "{%t}", l.Unterminated())
	}
}

type Character struct {
	l Literal
}

func (l Character) MissingCharacter() bool {
	return l.l.flags&literalFlagMissingChar != 0
}

func (l Character) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('+') {
		fmt.Fprintf(f, "%T{MissingCharacter:%t}", l, l.MissingCharacter())
	} else {
		fmt.Fprintf(f, "{%t}", l.MissingCharacter())
	}
}

type LiteralKind uint8

const (
	LiteralKindInteger LiteralKind = iota
	LiteralKindFloat
	LiteralKindString
	LiteralKindCharacter
)

type Base uint8

const (
	BaseBinary      Base = 2
	BaseOctal       Base = 8
	BaseDecimal     Base = 10
	BaseHexadecimal Base = 16
)

type TokenConsumer interface {
	ConsumeToken(Token)
}

func Tokenize(consumer TokenConsumer, s string) error {
	t := &tokenizer{consumer: consumer, s: s}
	for {
		token := t.Advance()
		if token.Kind() == kindNone {
			break
		}
		consumer.ConsumeToken(token)
	}
	return nil
}

type tokenizer struct {
	consumer     TokenConsumer
	s            string
	posWithinTok uint32
}

func (t *tokenizer) Advance() Token {
	token := t.nextToken()
	token.len = t.posWithinToken()
	t.resetPosWithinToken()
	return token
}

func (t *tokenizer) nextToken() Token {
	firstChar := t.bump()
	if firstChar == charEOF {
		return Token{kind: kindNone}
	}

	if t.isWhitespace(firstChar) {
		return t.whitespace()
	}

	if t.isSymbolStart(firstChar) {
		return t.symbol()
	}

	if t.isDecimal(firstChar) {
		return t.numeric(firstChar)
	}

	switch firstChar {
	case '\\':
		return t.character()
	case '"':
		return t.string()
	case ';':
		return t.lineComment()
	case '(':
		return Token{kind: KindOpenParen}
	case ')':
		return Token{kind: KindCloseParen}
	case '{':
		return Token{kind: KindOpenBrace}
	case '}':
		return Token{kind: KindCloseBrace}
	case '[':
		return Token{kind: KindOpenBracket}
	case ']':
		return Token{kind: KindCloseBracket}
	case '\'':
		return Token{kind: KindQuote}
	case '`':
		return Token{kind: KindBackquote}
	case '@':
		return Token{kind: KindDeref}
	case '^':
		return Token{kind: KindMetadata}
	case '#':
		return Token{kind: KindDispatch}
	case '~':
		return t.unquote()
	default:
		return Token{kind: KindInvalid}
	}
}

func (t *tokenizer) whitespace() Token {
	for t.isWhitespace(t.first()) {
		t.bump()
	}
	return Token{kind: KindWhitespace}
}

func (t *tokenizer) numeric(firstDigit rune) Token {
	var flags literalFlags

	if firstDigit == '0' {
		c := t.first()
		switch c {
		case 'b':
			t.bump()
			flags.setIf(!t.eatBinaryDigits(), literalFlagEmptyInt)
			return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindInteger, base: BaseBinary, flags: flags}}
		case 'o':
			t.bump()
			flags.setIf(!t.eatOctalDigits(), literalFlagEmptyInt)
			return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindInteger, base: BaseOctal, flags: flags}}
		case 'x':
			t.bump()
			flags.setIf(!t.eatHexDigits(), literalFlagEmptyInt)
			return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindInteger, base: BaseHexadecimal, flags: flags}}
		}
	}

	t.eatDecimalDigits()

	switch t.first() {
	case '.':
		if !t.isDecimalContinue(t.second()) {
			return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindInteger, base: BaseDecimal, flags: flags}}
		}

		t.bump()
		t.eatDecimalDigits()
		switch t.first() {
		case 'e', 'E':
			t.bump()
			flags.setIf(!t.eatFloatExponent(), literalFlagEmptyExponent)
		}
		return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindFloat, flags: flags}}
	case 'e', 'E':
		t.bump()
		flags.setIf(!t.eatFloatExponent(), literalFlagEmptyExponent)
		return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindFloat, flags: flags}}
	default:
		return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindInteger, base: BaseDecimal, flags: flags}}
	}
}

func (t *tokenizer) character() Token {
	var flags literalFlags
	c := t.bump()

	if t.isSymbolStart(c) {
		t.eatSymbol()
		return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindCharacter, flags: flags}}
	}

	flags.setIf(c == charEOF, literalFlagMissingChar)
	return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindCharacter, flags: flags}}
}

func (t *tokenizer) string() Token {
	var flags literalFlags

	for {
		c := t.bump()
		switch c {
		case charEOF:
			flags.setIf(true, literalFlagUnterminated)
			return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindString, flags: flags}}
		case '"':
			return Token{kind: KindLiteral, literal: Literal{kind: LiteralKindString, flags: flags}}
		case '\\':
			if t.first() == '"' || t.first() == '\\' {
				t.bump()
			}
		}
	}
}

func (t *tokenizer) lineComment() Token {
	for {
		c := t.first()
		if c == charEOF || c == '\n' {
			return Token{kind: KindLineComment}
		}
		t.bump()
	}
}

func (t *tokenizer) unquote() Token {
	if t.first() == '@' {
		t.bump()
		return Token{kind: KindUnquoteSplicing}
	}
	return Token{kind: KindUnquote}
}

func (t *tokenizer) isWhitespace(c rune) bool {
	switch c {
	case '\u0009', // \t
		'\u000A', // \n
		'\u000B', // vertical tab
		'\u000C', // form feed
		'\u000D', // \r
		'\u0020', // space
		'\u0085', // NEXT LINE from latin1
		'\u200E', // LEFT-TO-RIGHT MARK
		'\u200F', // RIGHT-TO-LEFT MARK
		'\u2028', // LINE SEPARATOR
		'\u2029': // PARAGRAPH SEPARATOR
		return true
	default:
		return false
	}
}

func (t *tokenizer) symbol() Token {
	t.eatSymbol()
	return Token{kind: KindSymbol}
}

func (t *tokenizer) eatSymbol() {
	for t.isSymbolContinue(t.first()) {
		t.bump()
	}
}

func (t *tokenizer) isSymbolStart(c rune) bool {
	if unicode.IsLetter(c) {
		return true
	}

	switch c {
	case '+',
		'-',
		'*',
		'/',
		'!',
		'_',
		'?',
		'<',
		'>',
		'=',
		'.',
		':':
		return true
	default:
		return false
	}
}

func (t *tokenizer) isSymbolContinue(c rune) bool {
	if t.isSymbolStart(c) {
		return true
	}

	if unicode.IsDigit(c) {
		return true
	}

	switch c {
	case '\'':
		return true
	default:
		return false
	}
}

func (t *tokenizer) eatFloatExponent() (found bool) {
	c := t.first()
	if c == '+' || c == '-' {
		t.bump()
	}
	return t.eatDecimalDigits()
}

func (t *tokenizer) eatBinaryDigits() (found bool) {
	for t.isBinaryContinue(t.first()) {
		found = true
		t.bump()
	}
	return found
}

func (t *tokenizer) isBinaryContinue(c rune) bool {
	return c == '0' || c == '1' || c == '_'
}

func (t *tokenizer) eatOctalDigits() (found bool) {
	for t.isOctalContinue(t.first()) {
		found = true
		t.bump()
	}
	return found
}

func (t *tokenizer) isOctalContinue(c rune) bool {
	return (c >= '0' && c <= '7') || c == '_'
}

func (t *tokenizer) eatHexDigits() (found bool) {
	for t.isHexContinue(t.first()) {
		found = true
		t.bump()
	}
	return found
}

func (t *tokenizer) isHexContinue(c rune) bool {
	return t.isDecimalContinue(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func (t *tokenizer) eatDecimalDigits() (found bool) {
	for t.isDecimalContinue(t.first()) {
		found = true
		t.bump()
	}
	return found
}

func (t *tokenizer) isDecimal(c rune) bool {
	return c >= '0' && c <= '9'
}

func (t *tokenizer) isDecimalContinue(c rune) bool {
	return t.isDecimal(c) || c == '_'
}

func (t *tokenizer) bump() rune {
	if len(t.s) == 0 {
		return charEOF
	}
	r, size := utf8.DecodeRuneInString(t.s)
	t.s = t.s[size:]
	t.posWithinTok += uint32(size)
	return r
}

func (t *tokenizer) first() rune {
	if len(t.s) == 0 {
		return charEOF
	}
	r, _ := utf8.DecodeRuneInString(t.s)
	return r
}

func (t *tokenizer) second() rune {
	if len(t.s) <= 1 {
		return charEOF
	}
	_, s := utf8.DecodeRuneInString(t.s)
	r, _ := utf8.DecodeRuneInString(t.s[s:])
	return r
}

func (t *tokenizer) posWithinToken() uint32 {
	return t.posWithinTok
}

func (t *tokenizer) resetPosWithinToken() {
	t.posWithinTok = 0
}

const charEOF rune = -1

type literalFlags uint8

const (
	literalFlagEmptyInt      literalFlags = 1 << 0
	literalFlagEmptyExponent literalFlags = 1 << 0
	literalFlagUnterminated  literalFlags = 1 << 0
	literalFlagMissingChar   literalFlags = 1 << 0
)

func (f *literalFlags) setIf(cond bool, flag literalFlags) {
	if cond {
		*f |= flag
	} else {
		*f &= ^flag
	}
}
