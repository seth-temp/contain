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


The 'alternate' branch is a POSIX-compliant shell script that utilizes
* chroot
* ulimit
* ip net-ns

# Main Branch
## Compiling
Probably requires go 1.6, since that is what it was written in  
>go get && go build -o lc

