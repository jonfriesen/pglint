# pglint
This library is used to validate PostgreSQL statements and running them against the ecpg library to determine if they are valid.

## requirements
- Go v1.11+
- [ecpg library](https://manpages.debian.org/experimental/libecpg-dev/ecpg.1.en.html)
    ```
    // Ubuntu
    sudo apt install libecpg-dev
    // RedHat & CentOS
    sudo yum install postgresql-devel
    // Arch Linux
    sudo pacman -S postgresql-libs
    ```

## usage

The generaly rule of thumb is no repsonse is good news.

### Library
Consume the library by creating a ECPG object.
```go
config := ecpg.Config{
    AddSemiColon:   *addSemicolons,
    TrimWhiteSpace: *trimWhitespace,
    QuestionMarks:  *allowQuestionMarks,
}

linter, err := ecpg.NewECPG(&config)
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

### CLI usage
`pglint-cli [options] <query>`
```
Usage of pglint-cli:
  -addSemicolons
        Automatically add semicolons to queries if not already present (will fail without a semicolon) (default true)
  -questionmarks
        Determines whether or not the to fill question mark placeholders with sample data for linting (default true)
  -trim
        Trim whitespace at head and tail of query (default true)
```

Eg. `pglint-cli -trim=false -questionmarks=false 'SELECT * FROM myTable;'`

### Docker

`docker run jonfriesen/pglint 'SELECT DATABASE_NAME FROM M_DATABASES;'`

#### Build yourself
Build: `docker build -t pglint .`
Run: `docker run pglint 'SELECT DATABASE_NAME FROM M_DATABASES;'`
