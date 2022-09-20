# flex-time

Flexible time parer.

flex-time tries to bring YYYY-MM-DDTHH:mm:ss.SSSZ format parsing to Go lang.

## Token Conversion Rule

- escape
  - escape single character by placing proceeding backward-slash (`\`).
  - escape bunch of characters by enclose with single quote.
- optional parts
  - make string inside `[]` as optional part.

Available tokens are shown in the table below:

| token     | go token           | description                     |
| --------- | ------------------ | ------------------------------- |
| []        | N/A                | escape as optional              |
| \\        | N/A                | escape one succeeding character |
| ''        | N/A                | escape quoted characters        |
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
- Then sort layouts by length in descending order.
  - so in above case, sorted like:
    - `2006-01-02T15:04:05.000Z07:00`,
    - `2006-01-02T15:04:05.000`,
    - `2006-01-02T15:04Z07:00`,
    - `2006-01-02T15Z07:00`,
    - `2006-01-02T15:04`,
    - `2006-01-02Z07:00`,
    - `2006-01-02T15`,
    - `2006-01-02`,
- Try parsing with layout one by one, longer to shorter.
- Return time.Time on first non-error.
- Return last error if all layouts fails.
