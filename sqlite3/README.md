# sqlite learn note

[official docs](https://sqlite.org/cli.html)

> [!NOTE]
> For learning sql, see <https://sqlbolt.com>

## OPEN FILE

```bash
sqlite3 <file name>
```

## DOT COMMANDS

```sql
.tables -- show tables
.save <file name> -- save changee
.open <file name> -- reopen or open a database
.quit -- exit the sqlite3
```

## CREATE

```sql
CREATE TABLE coffee_table (
    id INT PRIMARY KEY,
    name TEXT,
    region TEXT,
    roast TEXT
);
```

---

## INSERT

```sql
INSERT INTO coffee_table (id, name, region, roast)
VALUES (1, 'default route', 'ethiopia', 'light');
```

> [!NOTE]
> Use the ' instead of "

---

## SELECT

```sql
SELECT id, name FROM coffee_table;

-- select with condition
SELECT id, name FROM coffee_table 
WHERE id = 1;
```

---

## UPDATE

```sql
UPDATE coffee_table SET name = NULL where id = 1;
```

---

## DELETE

```sql
DELETE FROM coffee_table
WHERE id = 1;
```

## go-sqlite3 types

|go        | sqlite3           |
|----------|-------------------|
|nil       | null              |
|int       | integer           |
|int64     | integer           |
|float64   | float             |
|bool      | integer           |
|[]byte    | blob              |
|string    | text              |
|time.Time | timestamp/datetime|
