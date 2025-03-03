[ "{" "(" "[" ] @indent
[ "}" ")" "]" ] @dedent
(switch_statement
  body: (compound_statement "}" @match)
  (#set! indent.matchIndentOf parent.startPosition))

(["case" "default"] @match
  (#set! indent.matchIndentOf parent.parent.startPosition)
  (#set! indent.offsetIndent 1))

["case" "default"] @indent

(if_statement
  condition: (parenthesized_expression ")" @indent
  (#is? test.lastTextOnRow true)))

(if_statement
  condition: (_) @indent
  consequence: [
    (return_statement)
    (continue_statement)
    (break_statement)
  ] @dedent.next)
