package expression

import "fmt"

<<<<<<< HEAD
// ConditionMode will specify the types of the struct ConditionBuilder,
// representing the different types of Conditions (i.e. And, Or, Between, ...)
=======
// ConditionMode will specify the types of the struct ConditionBuilder
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
type ConditionMode int

const (
	// UnsetCond will catch errors if users make an empty ConditionBuilder
	UnsetCond ConditionMode = iota
	// EqualCond will represent the Equal Clause ConditionBuilder
	EqualCond
	// AndCond will represent the And Clause ConditionBuilder
	AndCond
<<<<<<< HEAD
	// BetweenCond will represent the Between ConditionBuilder
	BetweenCond
)

// ConditionBuilder will represent the ConditionExpressions in DynamoDB. It is
// composed of operands (OperandBuilder) and other conditions (ConditionBuilder)
// There are many different types of conditions, specified by ConditionMode.
// Users will be able to call the BuildExpression() method on a ConditionBuilder
// to create an Expression which can then be used for operation inputs into
// DynamoDB.
// More Information at: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html
=======
)

// ConditionBuilder will represent the ConditionExpressions
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
type ConditionBuilder struct {
	operandList   []OperandBuilder
	conditionList []ConditionBuilder
	Mode          ConditionMode
}

<<<<<<< HEAD
// Equal will create a ConditionBuilder with two OperandBuilders as children,
// representing the two operands that are being compared. The resulting
// ConditionBuilder can be used to build other Conditions or to create an
// Expression to be used in an operation input. This will be the function call.
//
// Example:
//
// condition := expression.Equal(expression.NewPath("foo"), expression.NewValue(5))
//
// anotherCondition := expression.Not(condition)	// Used in another condition
// expression, err := condition.BuildExpression()	// Used to make an Expression
=======
// Equal

// Equal will create a ConditionBuilder. This will be the function call
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
func Equal(left, right OperandBuilder) ConditionBuilder {
	return ConditionBuilder{
		operandList: []OperandBuilder{left, right},
		Mode:        EqualCond,
	}
}

// Equal will create a ConditionBuilder. This will be the method for PathBuilder
<<<<<<< HEAD
//
// Example:
//
// The following produce equivalent conditions:
// condition := expression.Equal(expression.NewPath("foo"), expression.NewValue(5))
// condition := expression.NewPath("foo").Equal(expression.NewValue(5))
=======
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
func (p PathBuilder) Equal(right OperandBuilder) ConditionBuilder {
	return Equal(p, right)
}

// Equal will create a ConditionBuilder. This will be the method for
// ValueBuilder
<<<<<<< HEAD
//
// Example:
//
// The following produce equivalent conditions:
// condition := expression.Equal(expression.NewValue(10), expression.NewValue(5))
// condition := expression.NewValue(10).Equal(expression.NewValue(5))
=======
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
func (v ValueBuilder) Equal(right OperandBuilder) ConditionBuilder {
	return Equal(v, right)
}

// Equal will create a ConditionBuilder. This will be the method for SizeBuilder
<<<<<<< HEAD
//
// Example:
//
// The following produce equivalent conditions:
// condition := expression.Equal(expression.NewPath("foo").Size(), expression.NewValue(5))
// condition := expression.NewPath("foo").Size().Equal(expression.NewValue(5))
=======
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
func (s SizeBuilder) Equal(right OperandBuilder) ConditionBuilder {
	return Equal(s, right)
}

<<<<<<< HEAD
// And will create a ConditionBuilder with more than two other Conditions as
// children, representing logical statements that will be logically ANDed
// together. The resulting ConditionBuilder can be used to build other
// Conditions or to create an Expression to be used in an operation input. This
// will be the function call.
//
// Example:
//
// condition1 := expression.Equal(expression.NewPath("foo"), expression.NewValue(5))
// condition2 := expression.Less(expression.NewPath("bar"), expression.NewValue(2010))
// condition3 := expression.NewPath("baz").Between(expression.NewValue(2), expression.NewValue(10))
// andCondition := expression.And(condition1, condition2, condition3)
//
// anotherCondition := expression.Not(andCondition)		// Used in another condition
// expression, err := andCondition.BuildExpression()	// Used to make an Expression
func And(cond ...ConditionBuilder) ConditionBuilder {
	return ConditionBuilder{
		conditionList: cond,
		Mode:          AndCond,
	}
}

// And will create a ConditionBuilder. This will be the method signature
//
// Example:
//
// The following produce equivalent conditions:
// condition := expression.And(condition1, condition2, condition3)
// condition := condition1.And(condition2, condition3)
func (cond ConditionBuilder) And(right ...ConditionBuilder) ConditionBuilder {
	right = append(right, cond)
	return And(right...)
}

// Between will create a ConditionBuilder with three operands as children, the
// first operand representing the operand being compared, the second operand
// representing the lower bound value of the first operand, and the third
// operand representing the upper bound value of the first operand. The
// resulting ConditionBuilder can be used to build other Conditions or to create
// an Expression to be used in an operation input. This will be the function
// call.
//
// Example:
//
// condition := expression.Between(expression.NewPath("foo"), expression.NewValue(2), expression.NewValue(6))
//
// anotherCondition := expression.Not(condition)	// Used in another condition
// expression, err := condition.BuildExpression()	// Used to make an Expression
func Between(ope, lower, upper OperandBuilder) ConditionBuilder {
	return ConditionBuilder{
		operandList: []OperandBuilder{ope, lower, upper},
		Mode:        BetweenCond,
	}
}

// Between will create a ConditionBuilder. This will be the method signature for
// PathBuilders.
//
// Example:
//
// The following produce equivalent conditions:
// condition := expression.Between(operand1, operand2, operand3)
// condition := operand1.Between(operand2, operand3)
func (p PathBuilder) Between(lower, upper OperandBuilder) ConditionBuilder {
	return Between(p, lower, upper)
}

// Between will create a ConditionBuilder. This will be the method signature for
// ValueBuilders.
//
// Example:
//
// The following produce equivalent conditions:
// condition := expression.Between(operand1, operand2, operand3)
// condition := operand1.Between(operand2, operand3)
func (v ValueBuilder) Between(lower, upper OperandBuilder) ConditionBuilder {
	return Between(v, lower, upper)
}

// Between will create a ConditionBuilder. This will be the method signature for
// SizeBuilders.
//
// Example:
//
// The following produce equivalent conditions:
// condition := expression.Between(operand1, operand2, operand3)
// condition := operand1.Between(operand2, operand3)
func (s SizeBuilder) Between(lower, upper OperandBuilder) ConditionBuilder {
	return Between(s, lower, upper)
}

// BuildExpression will take an ConditionBuilder as input and output an
// Expression which can be used in DynamoDB operational inputs (i.e.
// UpdateItemInput, DeleteItemInput, etc) In the future, the Expression struct
// can be used in some injection method into the input structs.
//
// Example:
//
// expr, err := someCondition.BuildExpression()
//
// deleteInput := dynamodb.DeleteItemInput{
// 	ConditionExpression:				aws.String(expr.Expression),
// 	ExpressionAttributeNames:		expr.Names,
// 	ExpressionAttributeValues:	expr.Values,
// 	Key: map[string]*dynamodb.AttributeValue{
// 		"PartitionKey": &dynamodb.AttributeValue{
// 			S: aws.String("SomeKey"),
// 		},
// 	},
// 	TableName: aws.String("SomeTable"),
// }
=======
// BuildExpression will take an ConditionBuilder as input and output an
// Expression
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
func (cond ConditionBuilder) BuildExpression() (Expression, error) {
	en, err := cond.buildCondition()
	if err != nil {
		return Expression{}, err
	}

	expr, err := en.buildExprNodes(&aliasList{})
	if err != nil {
		return Expression{}, err
	}

	return expr, nil
}

<<<<<<< HEAD
// buildCondition will build a tree structure of ExprNodes based on the tree
// structure of the input ConditionBuilder's child Conditions/Operands.
=======
// buildCondition will iterate over the tree of ConditionBuilders and
// OperandBuilders and build a tree of ExprNodes
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
func (cond ConditionBuilder) buildCondition() (ExprNode, error) {
	switch cond.Mode {
	case EqualCond:
		return compareBuildCondition(cond)
<<<<<<< HEAD
	case AndCond:
		return compoundBuildCondition(cond)
	case BetweenCond:
		return betweenBuildCondition(cond)
=======
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
	}
	return ExprNode{}, fmt.Errorf("No matching Mode to %v", cond.Mode)
}

// compareBuildCondition is the function to make ExprNodes from Compare
<<<<<<< HEAD
// ConditionBuilders. There will first be checks to make sure that the input
// ConditionBuilder has the correct format.
func compareBuildCondition(c ConditionBuilder) (ExprNode, error) {
	childNodes, err := buildChildNodes(c, 2, 0)
	if err != nil {
		return ExprNode{}, err
	}
	ret := ExprNode{
		children: childNodes,
	}

	// Create a string with special characters that can be substituted later: $c
	switch c.Mode {
	case EqualCond:
		ret.fmtExpr = "$c = $c"
	}

	return ret, nil
}

// compoundBuildCondition is the function to make ExprNodes from And/Or
// ConditionBuilders. There will first be checks to make sure that the input
// ConditionBuilder has the correct format.
func compoundBuildCondition(c ConditionBuilder) (ExprNode, error) {
	childNodes, err := buildChildNodes(c, 0, 2)
	if err != nil {
		return ExprNode{}, err
	}
	ret := ExprNode{
		children: childNodes,
	}

	// create a string with escaped characters to substitute them with proper
	// aliases during runtime
	for ind := range c.conditionList {
		ret.fmtExpr += "($c)"
		if ind != len(c.conditionList)-1 {
			switch c.Mode {
			case AndCond:
				ret.fmtExpr += " AND "
			}
		}
	}

	return ret, nil
}

// betweenBuildCondition is the function to make ExprNodes from Between
// ConditionBuilders. There will first be checks to make sure that the input
// ConditionBuilder has the correct format.
func betweenBuildCondition(c ConditionBuilder) (ExprNode, error) {
	childNodes, err := buildChildNodes(c, 3, 0)
	if err != nil {
		return ExprNode{}, err
	}
	ret := ExprNode{
		children: childNodes,
	}

	// Create a string with special characters that can be substituted later: $c
	ret.fmtExpr = "$c BETWEEN $c AND $c"

	return ret, nil
}

// buildChildNodes will check the format of the ConditionBuilder as well as
// create the list of the child ExprNodes. This avoids duplication of code
// amongst the various buildConditions.
func buildChildNodes(c ConditionBuilder, numOB, numCB int) ([]ExprNode, error) {
	if c.Mode == AndCond {
		if len(c.conditionList) < numCB {
			return []ExprNode{}, fmt.Errorf("Invalid ConditionBuilder. Expected at least %d Conditions", numCB)
		}
	} else {
		if len(c.conditionList) != numCB {
			return []ExprNode{}, fmt.Errorf("Invalid ConditionBuilder. Expected %d Conditions", numCB)
		}
	}

	// This check will be needed for In Condition, but not needed right now
	// if c.Mode == InCond {
	// 	if len(c.operandList) < numOB {
	// 		return []ExprNode{}, fmt.Errorf("Invalid ConditionBuilder. Expected at least %d Operands", numOB)
	// 	}
	// } else {
	if len(c.operandList) != numOB {
		return []ExprNode{}, fmt.Errorf("Invalid ConditionBuilder. Expected %d Operands", numOB)
	}
	//}

	var childNodes []ExprNode
	if len(c.operandList) == 0 {
		childNodes = make([]ExprNode, 0, len(c.conditionList))
		for _, cond := range c.conditionList {
			en, err := cond.buildCondition()
			if err != nil {
				return []ExprNode{}, err
			}
			childNodes = append(childNodes, en)
		}
	} else if len(c.conditionList) == 0 {
		childNodes = make([]ExprNode, 0, len(c.operandList))
		for _, ope := range c.operandList {
			en, err := ope.BuildOperand()
			if err != nil {
				return []ExprNode{}, err
			}
			childNodes = append(childNodes, en)
		}
	}

	return childNodes, nil
=======
// ConditionBuilders
func compareBuildCondition(c ConditionBuilder) (ExprNode, error) {
	if len(c.conditionList) != 0 {
		return ExprNode{}, fmt.Errorf("Invalid ConditionBuilder. Expected 0 ConditionBuilders")
	}

	if len(c.operandList) != 2 {
		return ExprNode{}, fmt.Errorf("Invalid ConditionBuilder. Expected 2 Operands")
	}

	operandExprNodes := make([]ExprNode, 0, len(c.operandList))
	for _, ope := range c.operandList {
		exprNodes, err := ope.BuildOperand()
		if err != nil {
			return ExprNode{}, err
		}
		operandExprNodes = append(operandExprNodes, exprNodes)
	}

	ret := ExprNode{
		children: operandExprNodes,
	}

	// Create a string with special characters that can be substituted later: $c
	switch c.Mode {
	case EqualCond:
		ret.fmtExpr = "$c = $c"
	}

	return ret, nil
>>>>>>> d753f586ade03ef8fa60c18037e2679ee424ee0b
}
