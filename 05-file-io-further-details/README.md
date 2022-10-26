# Chapter 05: File I/O: Further Details

## Exercise 5-2

### Question

Write a program that opens an existing file for writing with the `O_APPEND`
flag, and then seeks to the beginning of the file before writing some data.
Where does the data appear in the file? Why?

### Answer

The data appears at the end of the file. When a file is opened with `O_APPEND`,
seeking to the start of the file has no effect.

- [Go](go/cmd/seek_append)

## Exericse 5-3

### Question

This exercise is designed to demonstrate why the atomicity guaranteed by
opening a file with the `O_APPEND` flag is necessary. Write a program that
takes up to three command-line arguments:

```
$ atomic_append <filename> <num-bytes> [x]
```

This program should open the specified filename (creating it if necessary) and
append `num-bytes` bytes to the file by using `write()` to write a byte at a
time. By default, the program should open the file with the `O_APPEND` flag,
but if a third command-line argument (`x`) is supplied, the `O_APPEND` flag
should be omitted, and instead the program should perform an `lseek(fd, 0,
SEEK_END)` call before each `write()`.  Run two instances of this program at
the same time without the `x` argument to write 1 million bytes to the same
file:

```
$ atomic_append f1 1000000 & atomic_append f1 1000000
```

Repeat the same steps, writing to a different file, but this time specifying
the `x` argument:

```
$ atomic_append f2 1000000 x & atomic_append f2 1000000 x
```

List the sizes of the files `f1` and `f2` using `ls -l` and explain the
difference.

### Answer

`f1` and `f2` had different file sizes.

```
-rw-r--r--. 1 sudomateo sudomateo 2000000 Oct 19 22:46 f1
-rw-r--r--. 1 sudomateo sudomateo 1033636 Oct 19 22:46 f2
```

Since `f1` was opened with the `O_APPEND` flag and did not use `lseek()`, each
byte was written to the end of the file regardless which process wrote the
byte. Since `f2` was not opened with the `O_APPEND` flag and used `lseek()` to
seek to the end of the file before a write, each process maintained two
separate cursors for the file descriptor and wrote their bytes at their
respective cursor location. The `lseek()` and `write()` calls of each process
happened in a nondeterministic order which resulted in `f2` having a smaller
file size than `f1`.

- [Go](go/cmd/atomic_append)

## Exercise 5-4

### Question

Implement `dup()` and `dup2()` using `fcntl()` and, where necessary, `close()`.
(You may ignore the fact that `dup2()` and `fcntl()` return different `errno`
values for some error cases.) For `dup2()`, remember to handle the special case
where `oldfd` equals `newfd`. In this case, you should check whether `oldfd` is
valid, which can be done by, for example, checking if `fcntl(oldfd, F_GETFL)`
succeeds. If `oldfd` is not valid, then the function should return -1 with
`errrno` set to `EBADF`.

### Answer

- [Go](go/cmd/dup)

## Exercise 5-5

### Question

Write a program to verify that duplicated file descriptors share a file offset
value and open file status flags.

### Answer

- [Go](go/cmd/dup)

## Exercise 5-6

### Question

After each of the calls to `write()` in the following code, explain what the
content of the output file would be, and why:

```
fd1 = open(file, O_RDWR | O_CREAT | O_TRUNC, S_IRUSR | S_IWUSR);
fd2 = dup(fd1);
fd3 = open(file, O_RDWR);
write(fd1, "Hello,", 6);
write(fd2, " world", 6);
lseek(fd2, 0, SEEK_SET);
write(fd1, "HELLO,", 6);
write(fd3, "Gidday", 6);
```

### Answer

The output will be `Gidday world`.

- [Go](go/cmd/multiple_writes)
