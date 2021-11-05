- [Picker](#picker)
  - [Single Table Schema](#single-table-schema)
  - [Architecture](#architecture)

# Picker

## Single Table Schema
| Entity | PK        | SK               | GSI1PK    | GSI1SK           | type   |
| ------ | --------- | ---------------- | --------- | ---------------- | ------ |
| Room   | ROOM#UUID | ROOM#UUID        | USER#UUID | ROOM#UUID        | room   |
| Option | ROOM#UUID | ROOM_OPTION#UUID | USER#UUID | ROOM_OPTION#UUID | option |

## Architecture
<img src="./architecture.svg">