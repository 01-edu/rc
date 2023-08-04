### RC (Restrictions Checker)

This program analyzes a go source file and displays in standard output the imports, functions, slice types and loops used without authorization.

### Installation:

You will need go installed with version >=1.19

Clone this repository and simply run `go install .` inside this repository's root folder. After that you can run the program by invoking the `rc` command.

You can also use this program without installing by cloning this repository and run: `go run PATH/TO/REPO`.

### Usage:

This program expect at least one argument: a go file to be checked.

Just invoke the command by using: `rc [OPTIONS] FILE_TO_CHECK`.

The following options can be used:

- `-h` prints usage.
- `-config FILE` the path to a YAML configuration file. By default it will look for a `config.yml` in the directory you invoke the command from.

### By default:

- Allowed

  - All functions declared inside the source file are allowed
  - Slices of all types are allowed
  - Loops are allowed
  - Relative imports are allowed

- Disallowed

  - NO absolute imports are allowed
  - NO built-in functions are allowed.
  - NO casting is allowed

### Configuration file:

This program needs a configuration file to be ran.
The configuration file is expected to be a valid YAML, even if it's empty.

One or more of the following keys will be parsed from configuration and update the restrictions configuration list:

- `no-for:boolean` prohibits the use of `for` loops in the program or function, defaults to `false`.
- `no-slices:boolean` disallows the use of all slices types, defaults to `true`.
- `no-these-slices:strings array` disallows the slices of the specified types, defaults to `[]`.
- `no-relative-imports:boolean` disallows the use of relative imports, defaults to `false`.
- `no-lit:string` disallows character and string literals that match the regex pattern provided, defaults to an empty string.
- `allow-cast:boolean` allows casting to every built-in type, defaults to `true`.
- `allow-builtin:boolean` allows all builtin functions and casting to builtin types, defaults to `false`.
- `allowed-function:strings array` allowed imports and functions from a package, defaults to `[]`. The format is the following:
    - `<package>.*` for full imports (all functions from that package are allowed)
    - `<package>`.`<function>` for partial imports (only the function is allowed)
    - `<package>`.`<function>#amount` the function is only allowed to be used `amount` number of times
    - Examples:
        - `fmt.*` (all functions from `fmt` are allowed)
        - `github.com/01-edu/z01.PrintRune` (only `z01.PrintRune` is allowed)
        - `fmt.Println#2` (fmt.Println can only be used 2 times or less)
  - Allowed built-in functions
    - Use the name of the built-in function
    - It is also possible to limit the number of calls of a functions like with the imports using the '#' character

### Examples:

To allow the import of the whole `fmt` package, `z01.PrintRune` and the built-in functions `len` for the file `main.go`

Note: The imports must be written exactly the way they are written inside the source code, example:

  ```yaml
  allowed-functions: [fmt.*, github.com/01-edu/z01.PrintRune, len]
  ```

You can used also use the other format for arrays in yaml:
  ```yaml
  allowed-functions:
    - fmt.*
    - github.com/01-edu/z01.PrintRune
    - len
  ```


- Import "fmt" is allowed:

  ```yaml
  allowed-functions:
    - fmt.*
  ```

- Import "go/parser" is allowed:

  ```yaml
  allowed-functions:
    - go/parser.*
  ```

- Import "github.com/01-edu/z01" is allowed:

  ```yaml
  allowed-functions:
    - github.com/01-edu/z01.*
  ```

- Disallow litterals containing any digit (in strings AND numeric types):

  ```yaml
  no-lit: "[0-9]" # Better use quotes here to avoid yaml interpreting special characters
  ```

- Allow all type of casting:

  ```yaml
  allow-cast: true
  ```

- Disallow the use of the slices of type `string` and `int`

  ```yaml
  no-these-slices: [string, int]
  ```

- To allow casting to`rune` type only, add the type to the `allowed-functions` parameter:

  ```yaml
  allowed-functions:
    - rune
  ```

### How to read the error message

Let us look to an example snipped of code, let us imagine this code in a file called `main.go`:

```go
package main

import "fmt"

func main() {
	for _, v := range "abcdefghijklmnopqrstuvwxyz" {
		fmt.Println(v)
	}
	fmt.Println()
}
```

And the following config.yml file:

```yaml
allowed-functions: 
  - github.com/01-edu/z01.PrintRune
```

Now let us run the `rc` and understand the message

```console
$> rc main.go 
Parsing:
	Ok
Cheating:
	TYPE:             	NAME:      	LOCATION:
	illegal-import    	fmt        	main.go:3:8
	illegal-access    	fmt.Println	main.go:7:3
	illegal-access    	fmt.Println	main.go:10:2
	illegal-definition	main       	main.go:5:1
```

The important part is printed after the `Cheating` tag:

- The import of of the package `fmt` is not allowed
- In go the dot (.) is also known as the access operator for that reason the use of fmt.Println is shown as an illegal-access# Better use quotes here to avoid yaml interpreting special characters# Better use quotes here to avoid yaml interpreting special characters
- Finally the main function is shown as illegal-definition because the function is using disallowed functions that does not mean that the function can not be defined it just mean that the definition of the function must be changed to not use disallowed functions.
- Notice that the third column of the output with the tag "LOCATION:" show the location in the following way filepath:line:column
  This mean that you have to substitute the illegal function for ones that are allowed or write your own function with allowed functions

When `cheating` occurs, the exit code will be `127`, when the submitted code respect the restrictions, exit code will be `0`, and `1` will be returned in case an error occurs (no config file, invalid yaml, bad regex etc)
