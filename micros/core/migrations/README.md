Migration ceremony (lol):
1. Install [goose](https://github.com/pressly/goose)
2. `goose create <migration-step> sql`
3. Edit the resulting file, be aware that there are both (set) up and (tear) down part in
the migration sql script
4. To migrate up:
```sh
goose postgres "host=<host> port=<port> user=<user> password=<password> \
dbname=<dbname> sslmode=disable" up <migration-number> #default to the most recent
migration
```
5. To migrate down:
```sh
goose postgres "host=<host> port=<port> user=<user> password=<password> \
dbname=<dbname> sslmode=disable" down <migration-number> #default to current-1
``` 

`GOOSE_DBSTRING` could be set to the lengthy connection string above. Much recommended.
