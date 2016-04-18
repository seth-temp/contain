# Thought Process
Please note there are two branches in this repository: master and alternate
When I first saw this assignment, I thought a Docker wrapper might fit the 
requirements. After re-reading, I assumed this would be too easy to justify 
this project :).  

I think I may have misread the requirements. I missed the fact that 
assignment said  
"The contained program’s view of the file system should be limited to a
 specified directory"  


The main branch of this repo was written under the assumption that we were 
only looking for filesystem isolation. Therefore, I thought simply using 
LXC would be adequate.  


The 'alternate' branch is a POSIX-compliant shell script that utilizes
* chroot
* ulimit
* ip net-ns
