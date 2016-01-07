%% machine thrift;

package internal

import (
    "fmt"
    "io"
    "strconv"
    "strings"

    "github.com/uber/thriftrw-go/ast"
)

%%{
write data;

# Access state consistent across Lex() calls using the "lex" object.
access lex.;
variable p lex.p;
variable pe lex.pe;

}%%

type lexer struct {
    line int
    program *ast.Program

    err parseError
    parseFailed bool

    // Ragel:
    p, pe, cs, ts, te, act int
    data []byte
}


func newLexer(data []byte) *lexer {
    lex := &lexer{
        line: 1,
        err: newParseError(),
        parseFailed: false,
        data: data,
        p: 0,
        pe: len(data),
    }
    %% write init;
    return lex
}

func (lex *lexer) Lex(out *yySymType) int {
    eof := lex.pe
    tok := 0

    %%{
        ws = [ \t\r];

        # All uses of \n MUST use this instead if we want accurate line
        # number tracking.
        newline = '\n' >{ lex.line++ };

        __ = (ws | newline)*;

        # Comments
        line_comment = ('#'|'//') [^\n]*;
        multiline_comment = '/*' (newline | any)* :>> '*/';

        # Symbols are sent to the parser as-is.
        symbol = [\*=<>\(\)\{\},;:\[\]];

        # String literals.
        literal
            = ('"' ([^"\n\\] | '\\' any)* '"')
            | ("'" ([^'\n\\] | '\\' any)* "'")
            ;

        identifier = [a-zA-Z_] ([a-zA-Z0-9_] | '.' [a-zA-Z0-9_])*;

        integer = ('+' | '-')? digit+;
        hex_integer = '0x' xdigit+;

        double = integer '.' digit* ([Ee] integer)?;


        main := |*
            # A note about the scanner:
            #
            # Ragel will usually generate a scanner that will process all
            # available input in a single call. For goyacc, we want to advance
            # only to the next symbol and return that and any associated
            # information.
            #
            # So we use the special 'fbreak' statement available in action
            # blocks that consumes the token, saves the state, and breaks out
            # of the scanner. This allows the next call to the function to
            # pick up where the scanner left off.
            #
            # Because of this, we save all state for the scanner on the lexer
            # object.

            # Keywords
            'include'   __ => { tok =   INCLUDE; fbreak; };
            'namespace' __ => { tok = NAMESPACE; fbreak; };
            'void'      __ => { tok =      VOID; fbreak; };
            'bool'      __ => { tok =      BOOL; fbreak; };
            'byte'      __ => { tok =      BYTE; fbreak; };
            'i8'        __ => { tok =        I8; fbreak; };
            'i16'       __ => { tok =       I16; fbreak; };
            'i32'       __ => { tok =       I32; fbreak; };
            'i64'       __ => { tok =       I64; fbreak; };
            'double'    __ => { tok =    DOUBLE; fbreak; };
            'string'    __ => { tok =    STRING; fbreak; };
            'binary'    __ => { tok =    BINARY; fbreak; };
            'map'       __ => { tok =       MAP; fbreak; };
            'list'      __ => { tok =      LIST; fbreak; };
            'set'       __ => { tok =       SET; fbreak; };
            'oneway'    __ => { tok =    ONEWAY; fbreak; };
            'typedef'   __ => { tok =   TYPEDEF; fbreak; };
            'struct'    __ => { tok =    STRUCT; fbreak; };
            'union'     __ => { tok =     UNION; fbreak; };
            'exception' __ => { tok = EXCEPTION; fbreak; };
            'extends'   __ => { tok =   EXTENDS; fbreak; };
            'throws'    __ => { tok =    THROWS; fbreak; };
            'service'   __ => { tok =   SERVICE; fbreak; };
            'enum'      __ => { tok =      ENUM; fbreak; };
            'const'     __ => { tok =     CONST; fbreak; };
            'required'  __ => { tok =  REQUIRED; fbreak; };
            'optional'  __ => { tok =  OPTIONAL; fbreak; };
            'true'      __ => { tok =      TRUE; fbreak; };
            'false'     __ => { tok =     FALSE; fbreak; };

            symbol => {
                tok = int(lex.data[lex.ts])
                fbreak;
            };

            # Ignore comments and whitespace
            ws;
            newline;
            line_comment;
            multiline_comment;

            (integer | hex_integer) => {
                str := string(lex.data[lex.ts:lex.te])
                base := 10
                if len(str) > 2 && str[0:2] == "0x" {
                    // Hex constant
                    base = 16
                }

                if i64, err := strconv.ParseInt(str, base, 64); err != nil {
                    lex.Error(err.Error())
                } else {
                    out.i64 = i64
                    tok = INTCONSTANT
                }
                fbreak;
            };

            double => {
                str := string(lex.data[lex.ts:lex.te])
                if dub, err := strconv.ParseFloat(str, 64); err != nil {
                    lex.Error(err.Error())
                } else {
                    out.dub = dub
                    tok = DUBCONSTANT
                }
                fbreak;
            };

            literal => {
                bs := lex.data[lex.ts:lex.te]

                var str string
                var err error
                if len(bs) > 0 && bs[0] == '\'' {
                    str, err = UnquoteSingleQuoted(bs)
                } else {
                    str, err = strconv.Unquote(string(bs))
                }

                if err != nil {
                    lex.Error(err.Error())
                } else {
                    out.str = str
                    tok = LITERAL
                }

                fbreak;
            };

            identifier => {
                out.str = string(lex.data[lex.ts:lex.te])
                tok = IDENTIFIER
                fbreak;
            };
        *|;

        write exec;

    }%%

    if lex.cs == thrift_error {
        lex.Error(fmt.Sprintf("unknown token at index %d", lex.p))
    }
    return tok
}

func (lex *lexer) Error(e string) {
    lex.parseFailed = true
    lex.err.add(lex.line, e)
}
