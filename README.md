# GIM

Just like VIM but with "G" and written in GO.

# Getting started

```bash
go build ./cmd/gim
```
Now you can simply run the executable to use the editor.

```bash
./gim
```

# Usage
If you know VIM you're going to know GIM. \
It's the same but worse and with less functionallity😀

**Anyways here are the currently working commands:**
### Insertion
`i` - Enter insert _before_ cursor \
`I` - Enter insert _at beginning of line_ \
`a` - Enter insert _after_ cursor \
`A` - Enter insert _at end of line_ cursor 
### Movement 
`h` - Left \
`j` - Down \
`k` - Up \
`l` - Right 

### Command Mode
`:` - Enters command mode \
`:wq <filename>` - Write file and exit GIM \
`:q` - Exit GIM 