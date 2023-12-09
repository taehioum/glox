package parser

type Precedence int

/**
 * Defines the different precedence levels used by the infix parsers. These
 * determine how a series of infix expressions will be grouped. For example,
 * "a + b * c - d" will be parsed as "(a + (b * c)) - d" because "*" has higher
 * precedence than "+" and "-". Here, bigger numbers mean higher precedence.
 */
const (
	PrecedenceAssignment Precedence = iota + 1 // =
	PrecedenceOr                               // or
	PrecedenceAnd                              // and
	PrecedenceEquality                         // == !=
	PrecedenceComparison                       // > >= < <=
	PrecedenceTerm                             // - +
	PrecedenceFactor                           // / *
	PrecedenceUnary                            // ! -
	PrecedencePostfix                          // ++ --
)
