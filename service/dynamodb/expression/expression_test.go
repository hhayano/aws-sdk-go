// +build go1.7

package expression

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type exprErrorMode string

const (
	noExpressionError exprErrorMode = ""
	// invalidEscChar error will occer if the escape char '$' is either followed
	// by an unsupported character or if the escape char is the last character
	invalidEscChar = "invalid escape"
	// outOfRange error will occur if there are more escaped chars than there are
	// actual values to be aliased.
	outOfRange = "out of range"
	// nilAliasList error will occur if the aliasList passed in has not been
	// initialized
	nilAliasList = "AliasList is nil"
)

func TestBuildExpression(t *testing.T) {
	cases := []struct {
		name               string
		input              ExprNode
		expectedNames      map[string]*string
		expectedValues     map[string]*dynamodb.AttributeValue
		expectedExpression string
		err                exprErrorMode
	}{
		{
			name: "basic name",
			input: ExprNode{
				names:   []string{"foo"},
				fmtExpr: "$n",
			},

			expectedValues: map[string]*dynamodb.AttributeValue{},
			expectedNames: map[string]*string{
				"#0": aws.String("foo"),
			},
			expectedExpression: "#0",
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
			expectedNames: map[string]*string{},
			expectedValues: map[string]*dynamodb.AttributeValue{
				":0": &dynamodb.AttributeValue{
					N: aws.String("5"),
				},
			},
			expectedExpression: ":0",
		},
		{
			name: "nested path",
			input: ExprNode{
				names:   []string{"foo", "bar"},
				fmtExpr: "$n.$n",
			},

			expectedValues: map[string]*dynamodb.AttributeValue{},
			expectedNames: map[string]*string{
				"#0": aws.String("foo"),
				"#1": aws.String("bar"),
			},
			expectedExpression: "#0.#1",
		},
		{
			name: "nested path with index",
			input: ExprNode{
				names:   []string{"foo", "bar", "baz"},
				fmtExpr: "$n.$n[0].$n",
			},
			expectedValues: map[string]*dynamodb.AttributeValue{},
			expectedNames: map[string]*string{
				"#0": aws.String("foo"),
				"#1": aws.String("bar"),
				"#2": aws.String("baz"),
			},
			expectedExpression: "#0.#1[0].#2",
		},
		{
			name: "basic size",
			input: ExprNode{
				names:   []string{"foo"},
				fmtExpr: "size ($n)",
			},
			expectedValues: map[string]*dynamodb.AttributeValue{},
			expectedNames: map[string]*string{
				"#0": aws.String("foo"),
			},
			expectedExpression: "size (#0)",
		},
		{
			name: "duplicate path name",
			input: ExprNode{
				names:   []string{"foo", "foo"},
				fmtExpr: "$n.$n",
			},
			expectedValues: map[string]*dynamodb.AttributeValue{},
			expectedNames: map[string]*string{
				"#0": aws.String("foo"),
			},
			expectedExpression: "#0.#0",
		},
		{
			name: "equal expression",
			input: ExprNode{
				children: []ExprNode{
					ExprNode{
						names:   []string{"foo"},
						fmtExpr: "$n",
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

			expectedNames: map[string]*string{
				"#0": aws.String("foo"),
			},
			expectedValues: map[string]*dynamodb.AttributeValue{
				":0": &dynamodb.AttributeValue{
					N: aws.String("5"),
				},
			},
			expectedExpression: "#0 = :0",
		},
		{
			name: "missing char after $",
			input: ExprNode{
				names:   []string{"foo", "foo"},
				fmtExpr: "$n.$",
			},
			err: invalidEscChar,
		},
		{
			name: "names out of range",
			input: ExprNode{
				names:   []string{"foo"},
				fmtExpr: "$n.$n",
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
			err: invalidEscChar,
		},
		{
			name:               "empty ExprNode",
			input:              ExprNode{},
			expectedExpression: "",
		},
		{
			name:  "nil aliasList",
			input: ExprNode{},
			err:   nilAliasList,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var expr string
			var err error
			if c.err == nilAliasList {
				expr, err = c.input.BuildExpressionString(nil)
			} else {
				expr, err = c.input.BuildExpressionString(&AliasList{})
			}

			if c.err != noExpressionError {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expectedExpression, expr; !reflect.DeepEqual(a, e) {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}

func TestAliasValue(t *testing.T) {
	cases := []struct {
		name     string
		input    *AliasList
		expected string
		err      exprErrorMode
	}{
		{
			name:  "nil alias list",
			input: nil,
			err:   nilAliasList,
		},
		{
			name:     "first item",
			input:    &AliasList{},
			expected: ":0",
		},
		{
			name: "fifth item",
			input: &AliasList{
				valuesList: []dynamodb.AttributeValue{
					{},
					{},
					{},
					{},
				},
			},
			expected: ":4",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			str, err := c.input.aliasValue(dynamodb.AttributeValue{})

			if c.err != noExpressionError {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expected, str; e != a {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}

func TestAliasPath(t *testing.T) {
	cases := []struct {
		name      string
		inputList *AliasList
		inputName string
		expected  string
		err       exprErrorMode
	}{
		{
			name:      "nil alias list",
			inputList: nil,
			err:       nilAliasList,
		},
		{
			name:      "new unique item",
			inputList: &AliasList{},
			inputName: "foo",
			expected:  "#0",
		},
		{
			name: "duplicate item",
			inputList: &AliasList{
				namesList: []string{
					"foo",
					"bar",
				},
			},
			inputName: "foo",
			expected:  "#0",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			str, err := c.inputList.aliasPath(c.inputName)

			if c.err != noExpressionError {
				if err == nil {
					t.Errorf("expect error %q, got no error", c.err)
				} else {
					if e, a := string(c.err), err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %q error message to be in %q", e, a)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got unexpected Error %q", err)
				}

				if e, a := c.expected, str; e != a {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}
