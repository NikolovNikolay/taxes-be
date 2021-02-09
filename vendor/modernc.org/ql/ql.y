%{

//TODO Put your favorite license here
		
// yacc source generated by ebnf2y[1]
// at 2017-11-22 13:44:30.7008477 +0100 CET m=+0.004756809
//
//  $ ebnf2y -o ql.y -oe ql.ebnf -start StatementList -pkg ql -p _
//
// CAUTION: If this file is a Go source file (*.go), it was generated
// automatically by '$ go tool yacc' from a *.y file - DO NOT EDIT in that case!
// 
//   [1]: http://modernc.org/ebnf2y

package ql //TODO real package name

//TODO required only be the demo _dump function
import (
	"bytes"
	"fmt"
	"strings"

	"modernc.org/strutil"
)

%}

%union {
	item interface{} //TODO insert real field(s)
}

%token	_ANDAND
%token	_ANDNOT
%token	_EQ
%token	_FLOAT_LIT
%token	_GE
%token	_IDENTIFIER
%token	_IMAGINARY_LIT
%token	_INT_LIT
%token	_LE
%token	_LSH
%token	_NEQ
%token	_OROR
%token	_QL_PARAMETER
%token	_RSH
%token	_RUNE_LIT
%token	_STRING_LIT

%type	<item> 	/*TODO real type(s), if/where applicable */
	_ANDAND
	_ANDNOT
	_EQ
	_FLOAT_LIT
	_GE
	_IDENTIFIER
	_IMAGINARY_LIT
	_INT_LIT
	_LE
	_LSH
	_NEQ
	_OROR
	_QL_PARAMETER
	_RSH
	_RUNE_LIT
	_STRING_LIT

%token _ADD
%token _ALTER
%token _AND
%token _AS
%token _ASC
%token _BEGIN
%token _BETWEEN
%token _BIGINT
%token _BIGRAT
%token _BLOB
%token _BOOL
%token _BY
%token _BYTE
%token _COLUMN
%token _COMMIT
%token _COMPLEX128
%token _COMPLEX64
%token _CREATE
%token _DEFAULT
%token _DELETE
%token _DESC
%token _DISTINCT
%token _DROP
%token _DURATION
%token _EXISTS
%token _EXPLAIN
%token _FALSE
%token _FLOAT
%token _FLOAT32
%token _FLOAT64
%token _FROM
%token _FULL
%token _GROUPBY
%token _IF
%token _IN
%token _INDEX
%token _INSERT
%token _INT
%token _INT16
%token _INT32
%token _INT64
%token _INT8
%token _INTO
%token _IS
%token _JOIN
%token _LEFT
%token _LIKE
%token _LIMIT
%token _NOT
%token _NULL
%token _OFFSET
%token _ON
%token _OR
%token _ORDER
%token _OUTER
%token _RIGHT
%token _ROLLBACK
%token _RUNE
%token _SELECT
%token _SET
%token _STRING
%token _TABLE
%token _TIME
%token _TRANSACTION
%token _TRUE
%token _TRUNCATE
%token _UINT
%token _UINT16
%token _UINT32
%token _UINT64
%token _UINT8
%token _UNIQUE
%token _UPDATE
%token _VALUES
%token _WHERE

%type	<item> 	/*TODO real type(s), if/where applicable */
	AlterTableStmt
	AlterTableStmt1
	Assignment
	AssignmentList
	AssignmentList1
	AssignmentList2
	BeginTransactionStmt
	Call
	Call1
	Call11
	ColumnDef
	ColumnDef1
	ColumnDef11
	ColumnDef2
	ColumnName
	ColumnNameList
	ColumnNameList1
	ColumnNameList2
	CommitStmt
	Conversion
	CreateIndexStmt
	CreateIndexStmt1
	CreateIndexStmt2
	CreateTableStmt
	CreateTableStmt1
	CreateTableStmt2
	CreateTableStmt3
	DeleteFromStmt
	DeleteFromStmt1
	DropIndexStmt
	DropIndexStmt1
	DropTableStmt
	DropTableStmt1
	EmptyStmt
	ExplainStmt
	Expression
	Expression1
	Expression11
	ExpressionList
	ExpressionList1
	ExpressionList2
	Factor
	Factor1
	Factor11
	Factor2
	Field
	Field1
	FieldList
	FieldList1
	FieldList2
	GroupByClause
	Index
	IndexName
	InsertIntoStmt
	InsertIntoStmt1
	InsertIntoStmt2
	JoinClause
	JoinClause1
	JoinClause2
	Limit
	Literal
	Offset
	Operand
	OrderBy
	OrderBy1
	OrderBy11
	Predicate
	Predicate1
	Predicate11
	Predicate12
	Predicate121
	Predicate13
	PrimaryExpression
	PrimaryFactor
	PrimaryFactor1
	PrimaryFactor11
	PrimaryTerm
	PrimaryTerm1
	PrimaryTerm11
	QualifiedIdent
	QualifiedIdent1
	RecordSet
	RecordSet1
	RecordSet11
	RecordSet2
	RecordSetList
	RecordSetList1
	RecordSetList2
	RollbackStmt
	SelectStmt
	SelectStmt1
	SelectStmt2
	SelectStmt3
	SelectStmt4
	SelectStmt5
	SelectStmt6
	SelectStmt7
	SelectStmt8
	SelectStmt9
	Slice
	Slice1
	Slice2
	Start
	Statement
	StatementList
	StatementList1
	TableName
	Term
	Term1
	Term11
	TruncateTableStmt
	Type
	UnaryExpr
	UnaryExpr1
	UnaryExpr11
	UpdateStmt
	UpdateStmt1
	UpdateStmt2
	Values
	Values1
	Values2
	WhereClause

/*TODO %left, %right, ... declarations */

%start Start

%%

AlterTableStmt:
	_ALTER _TABLE TableName AlterTableStmt1
	{
		$$ = []AlterTableStmt{"ALTER", "TABLE", $3, $4} //TODO 1
	}

AlterTableStmt1:
	_ADD ColumnDef
	{
		$$ = []AlterTableStmt1{"ADD", $2} //TODO 2
	}
|	_DROP _COLUMN ColumnName
	{
		$$ = []AlterTableStmt1{"DROP", "COLUMN", $3} //TODO 3
	}

Assignment:
	ColumnName '=' Expression
	{
		$$ = []Assignment{$1, "=", $3} //TODO 4
	}

AssignmentList:
	Assignment AssignmentList1 AssignmentList2
	{
		$$ = []AssignmentList{$1, $2, $3} //TODO 5
	}

AssignmentList1:
	/* EMPTY */
	{
		$$ = []AssignmentList1(nil) //TODO 6
	}
|	AssignmentList1 ',' Assignment
	{
		$$ = append($1.([]AssignmentList1), ",", $3) //TODO 7
	}

AssignmentList2:
	/* EMPTY */
	{
		$$ = nil //TODO 8
	}
|	','
	{
		$$ = "," //TODO 9
	}

BeginTransactionStmt:
	_BEGIN _TRANSACTION
	{
		$$ = []BeginTransactionStmt{"BEGIN", "TRANSACTION"} //TODO 10
	}

Call:
	'(' Call1 ')'
	{
		$$ = []Call{"(", $2, ")"} //TODO 11
	}

Call1:
	/* EMPTY */
	{
		$$ = nil //TODO 12
	}
|	Call11
	{
		$$ = $1 //TODO 13
	}

Call11:
	'*'
	{
		$$ = "*" //TODO 14
	}
|	ExpressionList
	{
		$$ = $1 //TODO 15
	}

ColumnDef:
	ColumnName Type ColumnDef1 ColumnDef2
	{
		$$ = []ColumnDef{$1, $2, $3, $4} //TODO 16
	}

ColumnDef1:
	/* EMPTY */
	{
		$$ = nil //TODO 17
	}
|	ColumnDef11
	{
		$$ = $1 //TODO 18
	}

ColumnDef11:
	_NOT _NULL
	{
		$$ = []ColumnDef11{"NOT", "NULL"} //TODO 19
	}
|	Expression
	{
		$$ = $1 //TODO 20
	}

ColumnDef2:
	/* EMPTY */
	{
		$$ = nil //TODO 21
	}
|	_DEFAULT Expression
	{
		$$ = []ColumnDef2{"DEFAULT", $2} //TODO 22
	}

ColumnName:
	_IDENTIFIER
	{
		$$ = $1 //TODO 23
	}

ColumnNameList:
	ColumnName ColumnNameList1 ColumnNameList2
	{
		$$ = []ColumnNameList{$1, $2, $3} //TODO 24
	}

ColumnNameList1:
	/* EMPTY */
	{
		$$ = []ColumnNameList1(nil) //TODO 25
	}
|	ColumnNameList1 ',' ColumnName
	{
		$$ = append($1.([]ColumnNameList1), ",", $3) //TODO 26
	}

ColumnNameList2:
	/* EMPTY */
	{
		$$ = nil //TODO 27
	}
|	','
	{
		$$ = "," //TODO 28
	}

CommitStmt:
	_COMMIT
	{
		$$ = "COMMIT" //TODO 29
	}

Conversion:
	Type '(' Expression ')'
	{
		$$ = []Conversion{$1, "(", $3, ")"} //TODO 30
	}

CreateIndexStmt:
	_CREATE CreateIndexStmt1 _INDEX CreateIndexStmt2 IndexName _ON TableName '(' ExpressionList ')'
	{
		$$ = []CreateIndexStmt{"CREATE", $2, "INDEX", $4, $5, "ON", $7, "(", $9, ")"} //TODO 31
	}

CreateIndexStmt1:
	/* EMPTY */
	{
		$$ = nil //TODO 32
	}
|	_UNIQUE
	{
		$$ = "UNIQUE" //TODO 33
	}

CreateIndexStmt2:
	/* EMPTY */
	{
		$$ = nil //TODO 34
	}
|	_IF _NOT _EXISTS
	{
		$$ = []CreateIndexStmt2{"IF", "NOT", "EXISTS"} //TODO 35
	}

CreateTableStmt:
	_CREATE _TABLE CreateTableStmt1 TableName '(' ColumnDef CreateTableStmt2 CreateTableStmt3 ')'
	{
		$$ = []CreateTableStmt{"CREATE", "TABLE", $3, $4, "(", $6, $7, $8, ")"} //TODO 36
	}

CreateTableStmt1:
	/* EMPTY */
	{
		$$ = nil //TODO 37
	}
|	_IF _NOT _EXISTS
	{
		$$ = []CreateTableStmt1{"IF", "NOT", "EXISTS"} //TODO 38
	}

CreateTableStmt2:
	/* EMPTY */
	{
		$$ = []CreateTableStmt2(nil) //TODO 39
	}
|	CreateTableStmt2 ',' ColumnDef
	{
		$$ = append($1.([]CreateTableStmt2), ",", $3) //TODO 40
	}

CreateTableStmt3:
	/* EMPTY */
	{
		$$ = nil //TODO 41
	}
|	','
	{
		$$ = "," //TODO 42
	}

DeleteFromStmt:
	_DELETE _FROM TableName DeleteFromStmt1
	{
		$$ = []DeleteFromStmt{"DELETE", "FROM", $3, $4} //TODO 43
	}

DeleteFromStmt1:
	/* EMPTY */
	{
		$$ = nil //TODO 44
	}
|	WhereClause
	{
		$$ = $1 //TODO 45
	}

DropIndexStmt:
	_DROP _INDEX DropIndexStmt1 IndexName
	{
		$$ = []DropIndexStmt{"DROP", "INDEX", $3, $4} //TODO 46
	}

DropIndexStmt1:
	/* EMPTY */
	{
		$$ = nil //TODO 47
	}
|	_IF _EXISTS
	{
		$$ = []DropIndexStmt1{"IF", "EXISTS"} //TODO 48
	}

DropTableStmt:
	_DROP _TABLE DropTableStmt1 TableName
	{
		$$ = []DropTableStmt{"DROP", "TABLE", $3, $4} //TODO 49
	}

DropTableStmt1:
	/* EMPTY */
	{
		$$ = nil //TODO 50
	}
|	_IF _EXISTS
	{
		$$ = []DropTableStmt1{"IF", "EXISTS"} //TODO 51
	}

EmptyStmt:
	/* EMPTY */
	{
		$$ = nil //TODO 52
	}

ExplainStmt:
	_EXPLAIN Statement
	{
		$$ = []ExplainStmt{"EXPLAIN", $2} //TODO 53
	}

Expression:
	Term Expression1
	{
		$$ = []Expression{$1, $2} //TODO 54
	}

Expression1:
	/* EMPTY */
	{
		$$ = []Expression1(nil) //TODO 55
	}
|	Expression1 Expression11 Term
	{
		$$ = append($1.([]Expression1), $2, $3) //TODO 56
	}

Expression11:
	_OROR
	{
		$$ = $1 //TODO 57
	}
|	_OR
	{
		$$ = "OR" //TODO 58
	}

ExpressionList:
	Expression ExpressionList1 ExpressionList2
	{
		$$ = []ExpressionList{$1, $2, $3} //TODO 59
	}

ExpressionList1:
	/* EMPTY */
	{
		$$ = []ExpressionList1(nil) //TODO 60
	}
|	ExpressionList1 ',' Expression
	{
		$$ = append($1.([]ExpressionList1), ",", $3) //TODO 61
	}

ExpressionList2:
	/* EMPTY */
	{
		$$ = nil //TODO 62
	}
|	','
	{
		$$ = "," //TODO 63
	}

Factor:
	PrimaryFactor Factor1 Factor2
	{
		$$ = []Factor{$1, $2, $3} //TODO 64
	}

Factor1:
	/* EMPTY */
	{
		$$ = []Factor1(nil) //TODO 65
	}
|	Factor1 Factor11 PrimaryFactor
	{
		$$ = append($1.([]Factor1), $2, $3) //TODO 66
	}

Factor11:
	_GE
	{
		$$ = $1 //TODO 67
	}
|	'>'
	{
		$$ = ">" //TODO 68
	}
|	_LE
	{
		$$ = $1 //TODO 69
	}
|	'<'
	{
		$$ = "<" //TODO 70
	}
|	_NEQ
	{
		$$ = $1 //TODO 71
	}
|	_EQ
	{
		$$ = $1 //TODO 72
	}
|	_LIKE
	{
		$$ = "LIKE" //TODO 73
	}

Factor2:
	/* EMPTY */
	{
		$$ = nil //TODO 74
	}
|	Predicate
	{
		$$ = $1 //TODO 75
	}

Field:
	Expression Field1
	{
		$$ = []Field{$1, $2} //TODO 76
	}

Field1:
	/* EMPTY */
	{
		$$ = nil //TODO 77
	}
|	_AS _IDENTIFIER
	{
		$$ = []Field1{"AS", $2} //TODO 78
	}

FieldList:
	Field FieldList1 FieldList2
	{
		$$ = []FieldList{$1, $2, $3} //TODO 79
	}

FieldList1:
	/* EMPTY */
	{
		$$ = []FieldList1(nil) //TODO 80
	}
|	FieldList1 ',' Field
	{
		$$ = append($1.([]FieldList1), ",", $3) //TODO 81
	}

FieldList2:
	/* EMPTY */
	{
		$$ = nil //TODO 82
	}
|	','
	{
		$$ = "," //TODO 83
	}

GroupByClause:
	_GROUPBY ColumnNameList
	{
		$$ = []GroupByClause{"GROUP BY", $2} //TODO 84
	}

Index:
	'[' Expression ']'
	{
		$$ = []Index{"[", $2, "]"} //TODO 85
	}

IndexName:
	_IDENTIFIER
	{
		$$ = $1 //TODO 86
	}

InsertIntoStmt:
	_INSERT _INTO TableName InsertIntoStmt1 InsertIntoStmt2
	{
		$$ = []InsertIntoStmt{"INSERT", "INTO", $3, $4, $5} //TODO 87
	}

InsertIntoStmt1:
	/* EMPTY */
	{
		$$ = nil //TODO 88
	}
|	'(' ColumnNameList ')'
	{
		$$ = []InsertIntoStmt1{"(", $2, ")"} //TODO 89
	}

InsertIntoStmt2:
	Values
	{
		$$ = $1 //TODO 90
	}
|	SelectStmt
	{
		$$ = $1 //TODO 91
	}

JoinClause:
	JoinClause1 JoinClause2 _JOIN RecordSet _ON Expression
	{
		$$ = []JoinClause{$1, $2, "JOIN", $4, "ON", $6} //TODO 92
	}

JoinClause1:
	_LEFT
	{
		$$ = "LEFT" //TODO 93
	}
|	_RIGHT
	{
		$$ = "RIGHT" //TODO 94
	}
|	_FULL
	{
		$$ = "FULL" //TODO 95
	}

JoinClause2:
	/* EMPTY */
	{
		$$ = nil //TODO 96
	}
|	_OUTER
	{
		$$ = "OUTER" //TODO 97
	}

Limit:
	_LIMIT Expression
	{
		$$ = []Limit{"Limit", $2} //TODO 98
	}

Literal:
	_FALSE
	{
		$$ = "FALSE" //TODO 99
	}
|	_NULL
	{
		$$ = "NULL" //TODO 100
	}
|	_TRUE
	{
		$$ = "TRUE" //TODO 101
	}
|	_FLOAT_LIT
	{
		$$ = $1 //TODO 102
	}
|	_IMAGINARY_LIT
	{
		$$ = $1 //TODO 103
	}
|	_INT_LIT
	{
		$$ = $1 //TODO 104
	}
|	_RUNE_LIT
	{
		$$ = $1 //TODO 105
	}
|	_STRING_LIT
	{
		$$ = $1 //TODO 106
	}
|	_QL_PARAMETER
	{
		$$ = $1 //TODO 107
	}

Offset:
	_OFFSET Expression
	{
		$$ = []Offset{"OFFSET", $2} //TODO 108
	}

Operand:
	Literal
	{
		$$ = $1 //TODO 109
	}
|	QualifiedIdent
	{
		$$ = $1 //TODO 110
	}
|	'(' Expression ')'
	{
		$$ = []Operand{"(", $2, ")"} //TODO 111
	}

OrderBy:
	_ORDER _BY ExpressionList OrderBy1
	{
		$$ = []OrderBy{"ORDER", "BY", $3, $4} //TODO 112
	}

OrderBy1:
	/* EMPTY */
	{
		$$ = nil //TODO 113
	}
|	OrderBy11
	{
		$$ = $1 //TODO 114
	}

OrderBy11:
	_ASC
	{
		$$ = "ASC" //TODO 115
	}
|	_DESC
	{
		$$ = "DESC" //TODO 116
	}

Predicate:
	Predicate1
	{
		$$ = $1 //TODO 117
	}

Predicate1:
	Predicate11 Predicate12
	{
		$$ = []Predicate1{$1, $2} //TODO 118
	}
|	_IS Predicate13 _NULL
	{
		$$ = []Predicate1{"IS", $2, "NULL"} //TODO 119
	}

Predicate11:
	/* EMPTY */
	{
		$$ = nil //TODO 120
	}
|	_NOT
	{
		$$ = "NOT" //TODO 121
	}

Predicate12:
	_IN '(' ExpressionList ')'
	{
		$$ = []Predicate12{"IN", "(", $3, ")"} //TODO 122
	}
|	_IN '(' SelectStmt Predicate121 ')'
	{
		$$ = []Predicate12{"IN", "(", $3, $4, ")"} //TODO 123
	}
|	_BETWEEN PrimaryFactor _AND PrimaryFactor
	{
		$$ = []Predicate12{"BETWEEN", $2, "AND", $4} //TODO 124
	}

Predicate121:
	/* EMPTY */
	{
		$$ = nil //TODO 125
	}
|	';'
	{
		$$ = ";" //TODO 126
	}

Predicate13:
	/* EMPTY */
	{
		$$ = nil //TODO 127
	}
|	_NOT
	{
		$$ = "NOT" //TODO 128
	}

PrimaryExpression:
	Operand
	{
		$$ = $1 //TODO 129
	}
|	Conversion
	{
		$$ = $1 //TODO 130
	}
|	PrimaryExpression Index
	{
		$$ = []PrimaryExpression{$1, $2} //TODO 131
	}
|	PrimaryExpression Slice
	{
		$$ = []PrimaryExpression{$1, $2} //TODO 132
	}
|	PrimaryExpression Call
	{
		$$ = []PrimaryExpression{$1, $2} //TODO 133
	}

PrimaryFactor:
	PrimaryTerm PrimaryFactor1
	{
		$$ = []PrimaryFactor{$1, $2} //TODO 134
	}

PrimaryFactor1:
	/* EMPTY */
	{
		$$ = []PrimaryFactor1(nil) //TODO 135
	}
|	PrimaryFactor1 PrimaryFactor11 PrimaryTerm
	{
		$$ = append($1.([]PrimaryFactor1), $2, $3) //TODO 136
	}

PrimaryFactor11:
	'^'
	{
		$$ = "^" //TODO 137
	}
|	'|'
	{
		$$ = "|" //TODO 138
	}
|	'-'
	{
		$$ = "-" //TODO 139
	}
|	'+'
	{
		$$ = "+" //TODO 140
	}

PrimaryTerm:
	UnaryExpr PrimaryTerm1
	{
		$$ = []PrimaryTerm{$1, $2} //TODO 141
	}

PrimaryTerm1:
	/* EMPTY */
	{
		$$ = []PrimaryTerm1(nil) //TODO 142
	}
|	PrimaryTerm1 PrimaryTerm11 UnaryExpr
	{
		$$ = append($1.([]PrimaryTerm1), $2, $3) //TODO 143
	}

PrimaryTerm11:
	_ANDNOT
	{
		$$ = $1 //TODO 144
	}
|	'&'
	{
		$$ = "&" //TODO 145
	}
|	_LSH
	{
		$$ = $1 //TODO 146
	}
|	_RSH
	{
		$$ = $1 //TODO 147
	}
|	'%'
	{
		$$ = "%" //TODO 148
	}
|	'/'
	{
		$$ = "/" //TODO 149
	}
|	'*'
	{
		$$ = "*" //TODO 150
	}

QualifiedIdent:
	_IDENTIFIER QualifiedIdent1
	{
		$$ = []QualifiedIdent{$1, $2} //TODO 151
	}

QualifiedIdent1:
	/* EMPTY */
	{
		$$ = nil //TODO 152
	}
|	'.' _IDENTIFIER
	{
		$$ = []QualifiedIdent1{".", $2} //TODO 153
	}

RecordSet:
	RecordSet1 RecordSet2
	{
		$$ = []RecordSet{$1, $2} //TODO 154
	}

RecordSet1:
	TableName
	{
		$$ = $1 //TODO 155
	}
|	'(' SelectStmt RecordSet11 ')'
	{
		$$ = []RecordSet1{"(", $2, $3, ")"} //TODO 156
	}

RecordSet11:
	/* EMPTY */
	{
		$$ = nil //TODO 157
	}
|	';'
	{
		$$ = ";" //TODO 158
	}

RecordSet2:
	/* EMPTY */
	{
		$$ = nil //TODO 159
	}
|	_AS _IDENTIFIER
	{
		$$ = []RecordSet2{"AS", $2} //TODO 160
	}

RecordSetList:
	RecordSet RecordSetList1 RecordSetList2
	{
		$$ = []RecordSetList{$1, $2, $3} //TODO 161
	}

RecordSetList1:
	/* EMPTY */
	{
		$$ = []RecordSetList1(nil) //TODO 162
	}
|	RecordSetList1 ',' RecordSet
	{
		$$ = append($1.([]RecordSetList1), ",", $3) //TODO 163
	}

RecordSetList2:
	/* EMPTY */
	{
		$$ = nil //TODO 164
	}
|	','
	{
		$$ = "," //TODO 165
	}

RollbackStmt:
	_ROLLBACK
	{
		$$ = "ROLLBACK" //TODO 166
	}

SelectStmt:
	_SELECT SelectStmt1 SelectStmt2 SelectStmt3 SelectStmt4 SelectStmt5 SelectStmt6 SelectStmt7 SelectStmt8 SelectStmt9
	{
		$$ = []SelectStmt{"SELECT", $2, $3, $4, $5, $6, $7, $8, $9, $10} //TODO 167
	}

SelectStmt1:
	/* EMPTY */
	{
		$$ = nil //TODO 168
	}
|	_DISTINCT
	{
		$$ = "DISTINCT" //TODO 169
	}

SelectStmt2:
	'*'
	{
		$$ = "*" //TODO 170
	}
|	FieldList
	{
		$$ = $1 //TODO 171
	}

SelectStmt3:
	/* EMPTY */
	{
		$$ = nil //TODO 172
	}
|	_FROM RecordSetList
	{
		$$ = []SelectStmt3{"FROM", $2} //TODO 173
	}

SelectStmt4:
	/* EMPTY */
	{
		$$ = nil //TODO 174
	}
|	JoinClause
	{
		$$ = $1 //TODO 175
	}

SelectStmt5:
	/* EMPTY */
	{
		$$ = nil //TODO 176
	}
|	WhereClause
	{
		$$ = $1 //TODO 177
	}

SelectStmt6:
	/* EMPTY */
	{
		$$ = nil //TODO 178
	}
|	GroupByClause
	{
		$$ = $1 //TODO 179
	}

SelectStmt7:
	/* EMPTY */
	{
		$$ = nil //TODO 180
	}
|	OrderBy
	{
		$$ = $1 //TODO 181
	}

SelectStmt8:
	/* EMPTY */
	{
		$$ = nil //TODO 182
	}
|	Limit
	{
		$$ = $1 //TODO 183
	}

SelectStmt9:
	/* EMPTY */
	{
		$$ = nil //TODO 184
	}
|	Offset
	{
		$$ = $1 //TODO 185
	}

Slice:
	'[' Slice1 ':' Slice2 ']'
	{
		$$ = []Slice{"[", $2, ":", $4, "]"} //TODO 186
	}

Slice1:
	/* EMPTY */
	{
		$$ = nil //TODO 187
	}
|	Expression
	{
		$$ = $1 //TODO 188
	}

Slice2:
	/* EMPTY */
	{
		$$ = nil //TODO 189
	}
|	Expression
	{
		$$ = $1 //TODO 190
	}

Start:
	StatementList
	{
		_parserResult = $1 //TODO 191
	}

Statement:
	EmptyStmt
	{
		$$ = $1 //TODO 192
	}
|	AlterTableStmt
	{
		$$ = $1 //TODO 193
	}
|	BeginTransactionStmt
	{
		$$ = $1 //TODO 194
	}
|	CommitStmt
	{
		$$ = $1 //TODO 195
	}
|	CreateIndexStmt
	{
		$$ = $1 //TODO 196
	}
|	CreateTableStmt
	{
		$$ = $1 //TODO 197
	}
|	DeleteFromStmt
	{
		$$ = $1 //TODO 198
	}
|	DropIndexStmt
	{
		$$ = $1 //TODO 199
	}
|	DropTableStmt
	{
		$$ = $1 //TODO 200
	}
|	InsertIntoStmt
	{
		$$ = $1 //TODO 201
	}
|	RollbackStmt
	{
		$$ = $1 //TODO 202
	}
|	SelectStmt
	{
		$$ = $1 //TODO 203
	}
|	TruncateTableStmt
	{
		$$ = $1 //TODO 204
	}
|	UpdateStmt
	{
		$$ = $1 //TODO 205
	}
|	ExplainStmt
	{
		$$ = $1 //TODO 206
	}

StatementList:
	Statement StatementList1
	{
		$$ = []StatementList{$1, $2} //TODO 207
	}

StatementList1:
	/* EMPTY */
	{
		$$ = []StatementList1(nil) //TODO 208
	}
|	StatementList1 ';' Statement
	{
		$$ = append($1.([]StatementList1), ";", $3) //TODO 209
	}

TableName:
	_IDENTIFIER
	{
		$$ = $1 //TODO 210
	}

Term:
	Factor Term1
	{
		$$ = []Term{$1, $2} //TODO 211
	}

Term1:
	/* EMPTY */
	{
		$$ = []Term1(nil) //TODO 212
	}
|	Term1 Term11 Factor
	{
		$$ = append($1.([]Term1), $2, $3) //TODO 213
	}

Term11:
	_ANDAND
	{
		$$ = $1 //TODO 214
	}
|	_AND
	{
		$$ = "AND" //TODO 215
	}

TruncateTableStmt:
	_TRUNCATE _TABLE TableName
	{
		$$ = []TruncateTableStmt{"TRUNCATE", "TABLE", $3} //TODO 216
	}

Type:
	_BIGINT
	{
		$$ = "bigint" //TODO 217
	}
|	_BIGRAT
	{
		$$ = "bigrat" //TODO 218
	}
|	_BLOB
	{
		$$ = "blob" //TODO 219
	}
|	_BOOL
	{
		$$ = "bool" //TODO 220
	}
|	_BYTE
	{
		$$ = "byte" //TODO 221
	}
|	_COMPLEX128
	{
		$$ = "complex128" //TODO 222
	}
|	_COMPLEX64
	{
		$$ = "complex64" //TODO 223
	}
|	_DURATION
	{
		$$ = "duration" //TODO 224
	}
|	_FLOAT
	{
		$$ = "float" //TODO 225
	}
|	_FLOAT32
	{
		$$ = "float32" //TODO 226
	}
|	_FLOAT64
	{
		$$ = "float64" //TODO 227
	}
|	_INT
	{
		$$ = "int" //TODO 228
	}
|	_INT16
	{
		$$ = "int16" //TODO 229
	}
|	_INT32
	{
		$$ = "int32" //TODO 230
	}
|	_INT64
	{
		$$ = "int64" //TODO 231
	}
|	_INT8
	{
		$$ = "int8" //TODO 232
	}
|	_RUNE
	{
		$$ = "rune" //TODO 233
	}
|	_STRING
	{
		$$ = "string" //TODO 234
	}
|	_TIME
	{
		$$ = "time" //TODO 235
	}
|	_UINT
	{
		$$ = "uint" //TODO 236
	}
|	_UINT16
	{
		$$ = "uint16" //TODO 237
	}
|	_UINT32
	{
		$$ = "uint32" //TODO 238
	}
|	_UINT64
	{
		$$ = "uint64" //TODO 239
	}
|	_UINT8
	{
		$$ = "uint8" //TODO 240
	}

UnaryExpr:
	UnaryExpr1 PrimaryExpression
	{
		$$ = []UnaryExpr{$1, $2} //TODO 241
	}

UnaryExpr1:
	/* EMPTY */
	{
		$$ = nil //TODO 242
	}
|	UnaryExpr11
	{
		$$ = $1 //TODO 243
	}

UnaryExpr11:
	'^'
	{
		$$ = "^" //TODO 244
	}
|	'!'
	{
		$$ = "!" //TODO 245
	}
|	'-'
	{
		$$ = "-" //TODO 246
	}
|	'+'
	{
		$$ = "+" //TODO 247
	}

UpdateStmt:
	_UPDATE TableName UpdateStmt1 AssignmentList UpdateStmt2
	{
		$$ = []UpdateStmt{"UPDATE", $2, $3, $4, $5} //TODO 248
	}

UpdateStmt1:
	/* EMPTY */
	{
		$$ = nil //TODO 249
	}
|	_SET
	{
		$$ = "SET" //TODO 250
	}

UpdateStmt2:
	/* EMPTY */
	{
		$$ = nil //TODO 251
	}
|	WhereClause
	{
		$$ = $1 //TODO 252
	}

Values:
	_VALUES '(' ExpressionList ')' Values1 Values2
	{
		$$ = []Values{"VALUES", "(", $3, ")", $5, $6} //TODO 253
	}

Values1:
	/* EMPTY */
	{
		$$ = []Values1(nil) //TODO 254
	}
|	Values1 ',' '(' ExpressionList ')'
	{
		$$ = append($1.([]Values1), ",", "(", $4, ")") //TODO 255
	}

Values2:
	/* EMPTY */
	{
		$$ = nil //TODO 256
	}
|	','
	{
		$$ = "," //TODO 257
	}

WhereClause:
	_WHERE Expression
	{
		$$ = []WhereClause{"WHERE", $2} //TODO 258
	}
|	_WHERE _EXISTS '(' SelectStmt ')'
	{
		$$ = []WhereClause{"WHERE", "EXISTS", "(", $4, ")"} //TODO 259
	}
|	_WHERE _NOT _EXISTS '(' SelectStmt ')'
	{
		$$ = []WhereClause{"WHERE", "NOT", "EXISTS", "(", $5, ")"} //TODO 260
	}

%%

//TODO remove demo stuff below

var _parserResult interface{}

type (
	AlterTableStmt interface{}
	AlterTableStmt1 interface{}
	Assignment interface{}
	AssignmentList interface{}
	AssignmentList1 interface{}
	AssignmentList2 interface{}
	BeginTransactionStmt interface{}
	Call interface{}
	Call1 interface{}
	Call11 interface{}
	ColumnDef interface{}
	ColumnDef1 interface{}
	ColumnDef11 interface{}
	ColumnDef2 interface{}
	ColumnName interface{}
	ColumnNameList interface{}
	ColumnNameList1 interface{}
	ColumnNameList2 interface{}
	CommitStmt interface{}
	Conversion interface{}
	CreateIndexStmt interface{}
	CreateIndexStmt1 interface{}
	CreateIndexStmt2 interface{}
	CreateTableStmt interface{}
	CreateTableStmt1 interface{}
	CreateTableStmt2 interface{}
	CreateTableStmt3 interface{}
	DeleteFromStmt interface{}
	DeleteFromStmt1 interface{}
	DropIndexStmt interface{}
	DropIndexStmt1 interface{}
	DropTableStmt interface{}
	DropTableStmt1 interface{}
	EmptyStmt interface{}
	ExplainStmt interface{}
	Expression interface{}
	Expression1 interface{}
	Expression11 interface{}
	ExpressionList interface{}
	ExpressionList1 interface{}
	ExpressionList2 interface{}
	Factor interface{}
	Factor1 interface{}
	Factor11 interface{}
	Factor2 interface{}
	Field interface{}
	Field1 interface{}
	FieldList interface{}
	FieldList1 interface{}
	FieldList2 interface{}
	GroupByClause interface{}
	Index interface{}
	IndexName interface{}
	InsertIntoStmt interface{}
	InsertIntoStmt1 interface{}
	InsertIntoStmt2 interface{}
	JoinClause interface{}
	JoinClause1 interface{}
	JoinClause2 interface{}
	Limit interface{}
	Literal interface{}
	Offset interface{}
	Operand interface{}
	OrderBy interface{}
	OrderBy1 interface{}
	OrderBy11 interface{}
	Predicate interface{}
	Predicate1 interface{}
	Predicate11 interface{}
	Predicate12 interface{}
	Predicate121 interface{}
	Predicate13 interface{}
	PrimaryExpression interface{}
	PrimaryFactor interface{}
	PrimaryFactor1 interface{}
	PrimaryFactor11 interface{}
	PrimaryTerm interface{}
	PrimaryTerm1 interface{}
	PrimaryTerm11 interface{}
	QualifiedIdent interface{}
	QualifiedIdent1 interface{}
	RecordSet interface{}
	RecordSet1 interface{}
	RecordSet11 interface{}
	RecordSet2 interface{}
	RecordSetList interface{}
	RecordSetList1 interface{}
	RecordSetList2 interface{}
	RollbackStmt interface{}
	SelectStmt interface{}
	SelectStmt1 interface{}
	SelectStmt2 interface{}
	SelectStmt3 interface{}
	SelectStmt4 interface{}
	SelectStmt5 interface{}
	SelectStmt6 interface{}
	SelectStmt7 interface{}
	SelectStmt8 interface{}
	SelectStmt9 interface{}
	Slice interface{}
	Slice1 interface{}
	Slice2 interface{}
	Start interface{}
	Statement interface{}
	StatementList interface{}
	StatementList1 interface{}
	TableName interface{}
	Term interface{}
	Term1 interface{}
	Term11 interface{}
	TruncateTableStmt interface{}
	Type interface{}
	UnaryExpr interface{}
	UnaryExpr1 interface{}
	UnaryExpr11 interface{}
	UpdateStmt interface{}
	UpdateStmt1 interface{}
	UpdateStmt2 interface{}
	Values interface{}
	Values1 interface{}
	Values2 interface{}
	WhereClause interface{}
)
	
func _dump() {
	s := fmt.Sprintf("%#v", _parserResult)
	s = strings.Replace(s, "%", "%%", -1)
	s = strings.Replace(s, "{", "{%i\n", -1)
	s = strings.Replace(s, "}", "%u\n}", -1)
	s = strings.Replace(s, ", ", ",\n", -1)
	var buf bytes.Buffer
	strutil.IndentFormatter(&buf, ". ").Format(s)
	buf.WriteString("\n")
	a := strings.Split(buf.String(), "\n")
	for _, v := range a {
		if strings.HasSuffix(v, "(nil)") || strings.HasSuffix(v, "(nil),") {
			continue
		}
	
		fmt.Println(v)
	}
}

// End of demo stuff
