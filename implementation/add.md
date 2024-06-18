# add command

<b>description:</b>
used to add changes in the working directory to the staging area.

## Understanding the command
- It stores whatever changes there are in working directory, be it 
    - untracked files
    - modified files
    - deleted files 
    - new files 
    in index file
- It takes the data provided by status command about above types of file and updates the staging area accordingly 

## Think about it
### Initial steps
- Wait a second, but status command depends upon index file which is handled by add command
- And add command depends upon status of current working directory so that it can add the changes accordingly in the staging area
- Let's see how we can handle this deadlock
- If we observe the initial steps, initially there are only untracked files (after running init for the first time), so, for now we can build our status command so that it will only return untracked files and we will provide other data (deleted, modified and tracked files) once the cycle is up and running
- Now we will take the untracked files from status command and insert them using add command in the index file
- At this point the index file looks something like this

- By using above information we can create first stage of status command which checks for tracked files

### Dealing with modifications
- As discussed in status.md, to check if the file is modified or not we require its last modified time and its content (which we have in the index file)
- We can easily check the last modified time of files and compare them with the last modified time stored in index file
- But in case of no modification in content we will have to either decompress the file and check if the content is matching or not or
use the content of the file is see if the SHA1 has matches with the one that we have
- if it matches then the file is not modified and if it doesn't then the file is modified 

### Deletions
- If files entry exists in index file but file doesn't then it is deleted