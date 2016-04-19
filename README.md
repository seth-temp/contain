# Thought Process
When I first saw this assignment, I thought a Docker wrapper might fit the 
requirements. After re-reading, I assumed this would be too easy to justify 
this project.  

The project satisfies the following requirements
* chroot - By setting -p /path/to/lxc
* memory limit - By setting -m

I did not expect there to be library issues, but I discovered 
[bug #58](https://github.com/lxc/go-lxc/issues/58) in the go-lxc bindings.
It took me some time to question the library, as I incorrectly assumed 
I was doing something wrong (I am not as familiar with LXC as I am with Docker).  
After conversing with the developers in #lxc-dev on irc.freenode.net, I filed 
[pull request #59](https://github.com/lxc/go-lxc/pull/59). An unimpressive 3 
line fix, given the headache it caused me.

Also submitted a pull request for Execute method where stderr was being 
suppressed if the command exited nonzero 
[pull request #61](https://github.com/lxc/go-lxc/pull/61)

A much easier solution to this project would have been a POSIX script that uses
* chroot
* ulimit
* ip net-ns

# Main Branch
## Compiling
Probably requires go 1.6, since that is what it was written in, although 
with go's compat promise, you could build with 1.5 and GO15VENDOREXPERIMENT=1  
```sh
go get && go build -o contain
``
## Usage
```
./contain [command]
```

Use --help to print usage.
