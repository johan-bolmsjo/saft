# Saft

Saft is a data representation format similar to JSON with slight tweaks making
it more suitable to represent some types of data. It's the evolution of
[POT](https://github.com/johan-bolmsjo/pot) with the addition of raw strings,
comments and some questionable string escaping possibilities removed. Raw
strings are suitable to represent regexps in configuration files to avoid double
quoting.

This Go implementation focuses on parsing config files and does not provide a
stream parser. See the Go package documentation for API details. This document
mainly provide syntax details.

## Data Types

There are three distinct data types.

### List

A list begins with `[` and ends with `]`. Lists hold zero or more lists,
association lists and strings. Strings are separated by whitespace; it's
optional as element separator in other cases.

Examples:

    [] [[]] [[][]]
    [a a] [a "a a" `a a`]
    [a[a[a]]] [a [a [a]]]
    [a {a:a a:[a{a:a}]}]

### Association List

An association list is similar to a hash map but represented in list form. Each
list element is a pair of key and value. Ordering of pairs are maintained as
parsed and multiple pairs with identical key values are allowed.

An association list begins with `{` and ends with `}`. Keys and values are
separated by `:`. There may be optional whitespace after the colon but not
before it. Whitespace must follow the value of a pair unless the list is
terminated.

The key in a pair must be a string using either the symbol or interpreted syntax
form. The value may be of any type and string syntax form.

Examples:

    {}
    {a:b a:c}
    {a: {x:y} b: [i j k]}

### String

Everything boils down to a simple string of which there are three syntax forms.
It's up to the parsing application to interpret string data (e.g. convert to
numbers). Strings are encoded using UTF-8.

#### Symbol

The symbol form is unquoted and does not parse any escape sequences. It must not
start with character \` or `"` and is terminated by \\ \` " { } [ ] : and
whitespace.

#### Interpreted Quoted

Strings are quoted using character \". The form interprets escape codes \\n,
\\r, \\t with the usual meaning. Double quote must be escaped using \\" and
escape itself using \\\\. Interpreted strings may not span multiple lines
(contain \n or \r in their raw form).

#### Raw Quoted

Strings are quoted using character \`. The form does not interpret any escape
codes and can't represent the character \`. Raw strings may span multiple lines.

## Comments

Comments start with the character sequence `//` and stop at the end of the line. 

## Formal Grammar

TODO
