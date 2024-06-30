# concat
simple script to concatenate files

Usage
-----

### Example directory

```
$ ls
'file exists'  'file exists.001'  'file exists.002'  'file one.001'  'file one.002'  'file one.003'   go.mod   LICENSE   main.go  'missing part.002'  'missing part.003'   README.md
```

### Dryrun Concat files in working directory

```
$ concat -d
file 'file exists' already exists
wrong file name found. expected: 'missing part.001', found: 'missing part.002', file number: #1
```

### Concat files in random directory

```
$ concat ~/Projects/concat
file 'file exists' already exists
file created: file one
wrong file name found. expected: 'missing part.001', found: 'missing part.002', file number: #1
```