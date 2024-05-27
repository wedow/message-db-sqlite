# Message DB SQLite

A SQLite port of the [Message DB](https://github.com/message-db/message-db) event store.

## Installation

```bash
make && make install
```

## Usage

```sql
-- load the extension
select load_extension('message-db.so');

-- set up the schema
select initialize_message_db();

-- insert an event
select write_message('some ID', 'stream-1', 'event-1', 'type-1', 'data-1', 'metadata-1', -1);
```
