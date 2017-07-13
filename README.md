# fragment

Simple key-value-store server.

## commands

You can try with the telnet command.

The line starting from '>' is the response.
('> ' is not include)

The functions of the parameters are as follows.

| Name    | function                                               |
|---------|--------------------------------------------------------|
| key     | key                                                    |
| data    | data                                                   |
| flags   | 32bit integer                                          |
| exptime | expire time (second)                                   |
| bytes   | data size (byte)  Maximan size is 1000000 bytes(1MB)   |

### set

```
set [key] [flags] [exptime] [bytes]
[data]
> STORED
```

### get

```
get [key]
> VALUE [key] [flags] [bytes]
> [data]
> END
```

### replace

```
replace [key] [flags] [exptime] [bytes]
[data]
> STORED
```

### delete

```
delete [key]
> DELETED
```

## Lisence

MIT
