# Chapter 03: System Programming Concepts

## Exerise 3-1

### Question

When using the Linux-specific `reboot()` system call to reboot the system, the
second argument, `magic2`, must be specified as one of a set of magic numbers
(e.g. `LINUX_REBOOT_MAGIC2`). What is the significance of these numbers?
(Converting them to hexadecimal provides a clue.)

### Answer

The `magic` and `magic2` arguments are there to prevent erroneous calls to the
`reboot()` system call.

The manual page for `reboot()` lists the following accepted values for `magic`
and `magic2`:

```
LINUX_REBOOT_MAGIC1  = 0xFEE1DEAD (4276215469)
LINUX_REBOOT_MAGIC2  = 0x28121969 (672274793)
LINUX_REBOOT_MAGIC2A = 0x05121996 (85072278)
LINUX_REBOOT_MAGIC2B = 0x16041998 (369367448)
LINUX_REBOOT_MAGIC2C = 0x20112000 (537993216)
```

- `0xFEE1DEAD`: Feel dead. Humorous text since the system is going down for reboot.
- `0x28121969`: December 12, 1969. Linus Torvalds' birthday.
- `0x05121996`: December 5, 1996. Patricia Miranda's birthday. Linus Torvalds' daughter.
- `0x16041998`: April 16, 1998. Daniela Yolanda's birthday. Linus Torvalds' daughter.
- `0x20112000`: November 20, 2000. Celeste Amanda's birthday. Linus Torvalds' daughter.
