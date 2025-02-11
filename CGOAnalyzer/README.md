# CGOAnalyzer  

*CGOAnalyzer* is the tool introduced in our paper.

## Usage  

The tool allows specifying the repository storage path and the type of output information via command-line options.  

### Options:  
- `-h` : Display help information.  
- `-dumpdetail` : Write all information to a file, including paths. If not specified, path information will not be written.  
- `-showall` : Output all messages. If not specified, only errors and analyzed repository names will be displayed.  
- `-path` : Specify the directory to be analyzed (i.e., the directory containing all repositories to be analyzed). Default: `./`  
- `-core` : Number of core entries, referring to the top `x` (default: 10) most frequent entries.  
- `-top` : Number of top projects, referring to the top `x` (default: 5) repositories with the most core entries.  

### Example  

Suppose the repositories to be analyzed are stored in `/data/github_go/repos`, and you want to output detailed messages:  

```
$ go build ./
$ ./anatool -h
  Usage of ./anatool:
  -core int
    	number of core items (default 10)
  -dumppath
    	dump item path
  -path string
    	dir path of the repos (default "./")
  -showall
    	show output details
  -top int
    	top x repos having the core item (default 5)
$ ./anatool -path '/data/github_go/repos' -showall
```