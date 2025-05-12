# test data directory

This directoy contains some hand crafted database file used in tests.
Some of the test cases might require you to place a real database file
here, otherwise they would be skipped.

## File explanations


- `empty_scorelog.db`: An empty `scorelog.db` file
- `empty_songdata.db`: An empty `songdata.db` file
- `empty_scoredatalog.db`: An empty `scoredatalog.db` file

> These are real database files, you need to place your own here.
> These files must not be empty, since all empty file related test cases are already provided by `empty_*.db` files.
>
> Tests based on real files would skip if they are not provided.

- `scorelog.db`: A real `scorelog.db file, which is being gitignored
- `scoredatalog.db`: A real `scoredatalog.db file, which is being gitignored
- `songdata.db`: A real `songdata.db file, which is being gitignored
