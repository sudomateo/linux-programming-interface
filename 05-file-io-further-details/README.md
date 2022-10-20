# Chapter 05: File I/O: Further Details

## Exercise 5-2

### Question

Write a program that opens an existing file for writing with the `O_APPEND`
flag, and then seeks to the beginning of the file before writing some data.
Where does the data appear in the file? Why?

### Answer

The data appears at the end of the file. When a file is opened with `O_APPEND`,
seeking to the start of the file has no effect.

- [Go](go/cmd/append)
