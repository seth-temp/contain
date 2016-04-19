# Thought Process
When I first saw this assignment, I thought a Docker wrapper might fit the 
requirements. After re-reading, I assumed this would be too easy to justify 
this project.  

The project satisfies the following requirements
* chroot - By setting -p /path/to/lxc
* exit code matching - from LXC
* memory limit - By setting -m <bytes> !! UNTESTED

The project does not satisfy the following requirements due to deadline constraint
* Network Access - the container is given an IP address, but it is not assignable
* CPU affinity - I'm not even sure lxc-go can set this, would probably need
to wrap command with taskset(1)

I did not expect there to be library issues, but I discovered 
[bug #58](https://github.com/lxc/go-lxc/issues/58) in the go-lxc bindings.
It took me some time to question the library, as I incorrectly assumed 
I was doing something wrong (I am not as familiar with LXC as I am with Docker).  
After conversing with the developers in #lxc-dev on irc.freenode.net, I filed 
[pull request #59](https://github.com/lxc/go-lxc/pull/59). An unimpressive 3 
line fix, given the headache it caused me.

Also submitted a pull request for Execute method where stderr was being 
suppressed if the command exited nonzero.  
[pull request #61](https://github.com/lxc/go-lxc/pull/61)

# Application
## Compiling
Probably requires go 1.6, since that is what it was written in, although 
with go's compat promise, you could build with 1.5 and GO15VENDOREXPERIMENT=1  

Compiling may or may not require the lxc-dev (or equivalent) package 
```sh
go get && go build -o contain
```

## Requirements
The package lxc needs to be installed since the library 
[shells out to lxc-execute](https://github.com/lxc/go-lxc/blob/v2/container.go#L461-L473) 
as a fallback due to a bug.

## Usage
```
sudo ./contain <flags> [command]
IE
sudo ./contain --name foo --lxcpath /opt -- ls /bin
sudo ./contain -- du /
```

Use --help to print usage.

## Caveats
* If a container name (--name) is not specified, a brand new container is 
provisioned (with a random name)
* My machine generates initutils mount_fs warnings from the liblxc library

# Retrospective
Hindsight is 20/20, so they say. A much easier, deadline-meeting and 
feature complete solution to this project would have been a POSIX compliant
shell script that utilizes:
* chroot
* taskset
* ulimit
* ip net-ns

