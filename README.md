# flex-time

Flexible time parer.

flex-time tries to bring YYYY-MM-DDTHH:mm:ss.SSSZ format parsing to Go lang.

## Token Conversion Rule

- escape
  - escape single character by placing proceeding backward-slash (`\`).
- optional parts
  - make string inside `[]` as optional part.

Available tokens are shown in the table below:

| token     | go token           | description                     |
| --------- | ------------------ | ------------------------------- |
| []        | N/A                | escape as optional              |
| \\        | N/A                | escape one succeeding character |
| MMMM      | "January"          |                                 |
| MMM       | "Jan"              |                                 |
| M         | "1"                |                                 |
| MM        | "01"               |                                 |
| ww        | "Monday"           |                                 |
| w         | "Mon"              |                                 |
| d         | "2"                |                                 |
| dd        | "02"               |                                 |
| ddd       | "002"              |                                 |
| HH        | "15"               |                                 |
| h         | "3"                |                                 |
| hh        | "03"               |                                 |
| m         | "4"                |                                 |
| mm        | "04"               |                                 |
| s         | "5"                |                                 |
| ss        | "05"               |                                 |
| YYYY      | "2006"             |                                 |
| YY        | "06"               |                                 |
| A         | "PM"               |                                 |
| a         | "pm"               |                                 |
| MST       | "MST"              |                                 |
| ZZ        | "Z0700"            | prints Z for UTC                |
| Z070000   | "Z070000"          |                                 |
| Z07       | "Z07"              |                                 |
| Z         | "Z07:00"           | prints Z for UTC                |
| Z07:00:00 | "Z07:00:00"        |                                 |
| -0700     | "-0700"            | always numeric                  |
| -070000   | "-070000"          |                                 |
| -07       | "-07"              | always numeric                  |
| -07:00    | "-07:00"           | always numeric                  |
| -07:00:00 | "-07:00:00"        |                                 |
| .S[SS...] | ".0", ".00", ... , | trailing zeros included         |
| .0[00...] | ".0", ".00", ... , | trailing zeros included         |
| .9[99...] | ".9", ".99", ...,  | trailing zeros omitted          |

## Implementation

The implementation is pretty dumb.

- Fist of all, dump all patterns of given optional string. For example:
  - dump `YYYY-MM-DD[THH[:mm[:ss.SSS]]][z]` into:
    - `YYYY-MM-DDTHH:mm:ss.SSSZ`,
    - `YYYY-MM-DDTHH:mm:ss.SSS`,
    - `YYYY-MM-DDTHH:mmZ`,
    - `YYYY-MM-DDTHH:mm`,
    - `YYYY-MM-DDTHHZ`,
    - `YYYY-MM-DDTHH`,
    - `YYYY-MM-DDZ`,
    - `YYYY-MM-DD`,
- Convert all the `YYYY` `MM` things into golang time layout tokens, `2006` and `01`.
- Then count length of input date-time string. Use nearest longer layout first.
- If time.(`Parse`|`ParseInLocation`) with one layout fails, try one shorter format.
- If no shorter
