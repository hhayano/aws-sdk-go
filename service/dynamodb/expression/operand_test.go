// +build go1.8

package expression

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// opeErrorMode will help with error cases and checking error types
type opeErrorMode int

const (
	emptyPath opeErrorMode = iota + 1
	pathIndex
	escChar
	outOfRange
	nilAliasList
)

func (oem opeErrorMode) String() string {
	switch oem {
	case emptyPath:
		return "path is empty"
	case pathIndex:
		return "invalid path index"
	case escChar:
		return "invalid escape"
	case outOfRange:
		return "out of range"
	case nilAliasList:
		return "aliasList is nil"
	default:
		return "no matching opeErrorMode"
	}
}

func TestBuildOperand(t *testing.T) {
	cases := []struct {
		name     string
		input    OperandBuilder
		expected ExprNode
		err      opeErrorMode
	}{
		{
			name:  "basic path",
			input: Path("foo"),
			expected: ExprNode{
				names:   []string{"foo"},
				fmtExpr: "$p",
			},
		},
		{
			name:  "duplicate path name",
			input: Path("foo.foo"),
			expected: ExprNode{
				names:   []string{"foo", "foo"},
				fmtExpr: "$p.$p",
			},
		},
		{
			name:  "basic value",
			input: Value(5),
			expected: ExprNode{
				values: []dynamodb.AttributeValue{
					dynamodb.AttributeValue{
						N: aws.String("5"),
					},
				},
				fmtExpr: "$v",
			},
		},
		{
			name:  "nested path",
			input: Path("foo.bar"),
			expected: ExprNode{
				names:   []string{"foo", "bar"},
				fmtExpr: "$p.$p",
			},
		},
		{
			name:  "nested path with index",
			input: Path("foo.bar[0].baz"),
			expected: ExprNode{
				names:   []string{"foo", "bar", "baz"},
				fmtExpr: "$p.$p[0].$p",
			},
		},
		{
			name:  "basic size",
			input: Path("foo").Size(),
			expected: ExprNode{
				names:   []string{"foo"},
				fmtExpr: "size ($p)",
			},
		},
		{
			name:     "empty path error",
			input:    Path(""),
			expected: ExprNode{},
			err:      emptyPath,
		},
		{
			name:     "invalid path",
			input:    Path("foo..bar"),
			expected: ExprNode{},
			err:      emptyPath,
		},
		{
			name:     "invalid index",
			input:    Path("[foo]"),
			expected: ExprNode{},
			err:      pathIndex,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			en, err := c.input.BuildOperand()

			if c.err != 0 {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := c.err.String(), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expected, en; !reflect.DeepEqual(a, e) {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}

func TestBuildExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    ExprNode
		expected Expression
		err      opeErrorMode
	}{
		{
			name: "basic path",
			input: ExprNode{
				names:   []string{"foo"},
				fmtExpr: "$p",
			},
			expected: Expression{
				Names: map[string]*string{
					"#0": aws.String("foo"),
				},
				Expression: "#0",
			},
		},
		{
			name: "basic value",
			input: ExprNode{
				values: []dynamodb.AttributeValue{
					dynamodb.AttributeValue{
						N: aws.String("5"),
					},
				},
				fmtExpr: "$v",
			},
			expected: Expression{
				Values: map[string]*dynamodb.AttributeValue{
					":0": &dynamodb.AttributeValue{
						N: aws.String("5"),
					},
				},
				Expression: ":0",
			},
		},
		{
			name: "nested path",
			input: ExprNode{
				names:   []string{"foo", "bar"},
				fmtExpr: "$p.$p",
			},
			expected: Expression{
				Names: map[string]*string{
					"#0": aws.String("foo"),
					"#1": aws.String("bar"),
				},
				Expression: "#0.#1",
			},
		},
		{
			name: "nested path with index",
			input: ExprNode{
				names:   []string{"foo", "bar", "baz"},
				fmtExpr: "$p.$p[0].$p",
			},
			expected: Expression{
				Names: map[string]*string{
					"#0": aws.String("foo"),
					"#1": aws.String("bar"),
					"#2": aws.String("baz"),
				},
				Expression: "#0.#1[0].#2",
			},
		},
		{
			name: "basic size",
			input: ExprNode{
				names:   []string{"foo"},
				fmtExpr: "size ($p)",
			},
			expected: Expression{
				Names: map[string]*string{
					"#0": aws.String("foo"),
				},
				Expression: "size (#0)",
			},
		},
		{
			name: "duplicate path name",
			input: ExprNode{
				names:   []string{"foo", "foo"},
				fmtExpr: "$p.$p",
			},
			expected: Expression{
				Names: map[string]*string{
					"#0": aws.String("foo"),
				},
				Expression: "#0.#0",
			},
		},
		{
			name: "equal expression",
			input: ExprNode{
				children: []ExprNode{
					ExprNode{
						names:   []string{"foo"},
						fmtExpr: "$p",
					},
					ExprNode{
						values: []dynamodb.AttributeValue{
							dynamodb.AttributeValue{
								N: aws.String("5"),
							},
						},
						fmtExpr: "$v",
					},
				},
				fmtExpr: "$c = $c",
			},
			expected: Expression{
				Names: map[string]*string{
					"#0": aws.String("foo"),
				},
				Values: map[string]*dynamodb.AttributeValue{
					":0": &dynamodb.AttributeValue{
						N: aws.String("5"),
					},
				},
				Expression: "#0 = :0",
			},
		},
		{
			name: "missing char after $",
			input: ExprNode{
				names:   []string{"foo", "foo"},
				fmtExpr: "$p.$",
			},
			err: escChar,
		},
		{
			name: "names out of range",
			input: ExprNode{
				names:   []string{"foo"},
				fmtExpr: "$p.$p",
			},
			err: outOfRange,
		},
		{
			name: "values out of range",
			input: ExprNode{
				fmtExpr: "$v",
			},
			err: outOfRange,
		},
		{
			name: "children out of range",
			input: ExprNode{
				fmtExpr: "$c",
			},
			err: outOfRange,
		},
		{
			name: "invalid escape char",
			input: ExprNode{
				fmtExpr: "$!",
			},
			err: escChar,
		},
		{
			name:     "empty ExprNode",
			input:    ExprNode{},
			expected: Expression{},
		},
		{
			name:     "nil aliasList",
			input:    ExprNode{},
			expected: Expression{},
			err:      nilAliasList,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var expr Expression
			var err error
			if c.err == nilAliasList {
				expr, err = c.input.buildExprNodes(nil)
			} else {
				expr, err = c.input.buildExprNodes(&aliasList{})
			}

			if c.err != 0 {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := c.err.String(), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expected, expr; !reflect.DeepEqual(a, e) {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}
