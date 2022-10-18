For syscall, we get 3 returned values (r1, r2 uintptr, lastErr error)
r1 is used and compared with constants that we defined (such as internal/win/errors.go)
r2 is rarely/never used
lastErr is used and compared with errors offered in golang.org/x/sys/windows
___
At Microsoft, some variables are named after their types example:
pSomething is probably a pointer to Something
pllThreshold is a pointer to llThreshold
However this is not true in this code, variable names means nothing but a name
There are many cases where pllThreshold is just a regular value and it's pointer is defined by &pllThreshold
such variables naming was only so the coding is easier to write given that a lot of a code was a copy/paste from Microsoft Docs
___
Authentication is based on the user SID (auth package)
