# Thought Process
When I first saw this assignment, I thought a Docker wrapper might fit the 
requirements. After re-reading, I assumed this would be too easy to justify 
this project.  

The project does satisfy these following requirements:
* chroot - By setting -p /path/to/lxc
* exit code / signal forwarding - from lxc-init
* memory limit - By setting -m <bytes> !! UNTESTED

The project does not satisfy these following requirements due to deadline constraint
* Network Access - the container is given an IP address, but it is not assignable
* CPU affinity - I'm not even sure lxc-go can set this, would probably need
to wrap command with taskset(1)

I did not expect there to be library issues, but I discovered 
[issue #58](https://github.com/lxc/go-lxc/issues/58) in the go-lxc bindings.
It took me some time to question the library, as I incorrectly assumed 
I was doing something wrong (I am not as familiar with LXC semantics and commands 
as I am with Docker).  
After conversing with the developers in #lxc-dev on irc.freenode.net, I filed 
[pull request #59](https://github.com/lxc/go-lxc/pull/59). An unimpressive 3 
line fix, given the headache it caused me.

Also submitted a pull request for Execute method where stderr was being 
suppressed if the command exited nonzero.  
[pull request #61](https://github.com/lxc/go-lxc/pull/61)

# Application
## Requirements
If you are on a Debian-based distro, please ensure you are using the backports 
sources.list as it will otherwise default to the stable / incompatible LXC 1.0.8!
The following packages and specified versions are **required**
* lxc **version 1.1.5+** for the binary, since it [shells out to lxc-execute](https://github.com/lxc/go-lxc/blob/v2/container.go#L461-L473) 
* lxc-dev **version 1.1.5+** for compiling
* pkg-config for compiling the LXC C extensions

## Compiling
Requires go 1.6, since that is what it was written in. It may compile 
with go 1.5 if you export GO15VENDOREXPERIMENT=1 and believe in compatibility 
promises :)

Compiling may or may not require the lxc-dev (or equivalent) package 
```sh
go get && go build -o contain
```

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

