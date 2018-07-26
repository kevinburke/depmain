# depmain

depmain is designed to make it easy to test your program's main function. Create
a second main() function that accepts a `*depmain.Ext` and returns an integer.
Then call it from main like this:

```go
os.Exit(_main(depmain.New()))
```

Then in tests, you can replace `*depmain.Ext.Stdout` with a bytes.Buffer,
configure the environment/process args as you see fit, etc.

## Program Changes

You need to change what you write in your program to make sure you are using the
provided environment and not reading/writing the actual external environment.

#### flag.Parse

Replace instances of `flag.Parse` with `flag.CommandLine.Parse(ext.Args)`

```go
func _main(ext *depmain.Ext) int {
    if err := flag.CommandLine.Parse(ext.Args); err != nil {
        flag.Usage()
        return 2
    }
}
```

Note you'll also want to set `flag.Usage` to write to `ext.Stderr` not
`os.Stderr`.

#### os.Stdout, os.Stderr

Replace instances of `fmt.Println("hello")` with `fmt.Fprintln(ext.Stdout,
"hello")`.

#### os.Environ

Use `ext.Getenv` and `ext.LookupEnv` instead of `os.Getenv` and `os.LookupEnv`.

```go
func _main(ext *depmain.Ext) int {
    fmt.Println(ext.Getenv("TZ"))
    if u, found := ext.LookupEnv("DATABASE_URL"); found {
        connect(u)
    }
}
```
