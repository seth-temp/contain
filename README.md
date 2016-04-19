# Thought Process
Please note there are two branches in this repository: master and alternate.
When I first saw this assignment, I thought a Docker wrapper might fit the 
requirements. After re-reading, I assumed this would be too easy to justify 
this project :).  

I think I may have misread the requirements. I missed the fact that 
assignment said  
"The contained program’s view of the file system should be limited to a
 specified directory"  

The main branch of this repo was written under the assumption that we were 
only looking for filesystem isolation and not a chroot-anywhere feature. 
Therefore, I thought simply using LXC would be adequate. If LXC is fed a 
template, it acquires a filesystem of the distro specified.

I did not expect there to be library issues, but I discovered 
[bug #58](https://github.com/lxc/go-lxc/issues/58) in the go-lxc bindings.
It took me some time to question the library, as I incorrectly assumed 
I was doing something wrong (I am not as familiar with LXC as I am with Docker).  
After conversing with the developers in #lxc-dev on irc.freenode.net, I filed 
[pull request #59](https://github.com/lxc/go-lxc/pull/59). An unimpressive 3 
line fix, given the headache it caused me.


The 'alternate' branch is a POSIX-compliant shell script that utilizes
* chroot
* ulimit
* ip net-ns

# Main Branch
## Compiling
Probably requires go 1.6, since that is what it was written in, although 
with go's compat promise, you could build with 1.5 and GO15VENDOREXPERIMENT=1  
```sh
go get && go build -o contain
```

## Usage
```
./contain [command]
```

Use --help to print usage.
