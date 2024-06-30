# concat
simple script to concatenate files

Usage
-----

### Example directory

```
$ ls
'file exists.txt'  'file exists.txt.001'  'file exists.txt.002'  'file one.001'  'file one.002'  'file one.003'   go.mod   LICENSE   main.go  'missing part.002'  'missing part.003'   README.md
```

### Dry-Run Concat files in working directory

```
$ concat -d
file 'file exists' already exists
wrong file name found. expected: 'missing part.001', found: 'missing part.002', file number: #1
```

### Concat files in provided directory

```
$ concat ~/Projects/concat
file 'file exists.txt' already exists
file created: file one
wrong file name found. expected: 'missing part.001', found: 'missing part.002', file number: #1
```

### Concat files in forced mode will create unique file names and delete the used parts

```
$ concat ~/Projects/concat -f
file 'file exists.txt' already exists
file created: file exists-001.txt
file removed: file exists.txt.001
file removed: file exists.txt.002
file created: file one
file removed: file one.001
file removed: file one.002
file removed: file one.003
wrong file name found. expected: 'missing part.001', found: 'missing part.002', file number: #1
```
