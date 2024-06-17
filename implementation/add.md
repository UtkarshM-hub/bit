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
- Wait a second, but status command depends upon index file which is handled by add command
- And add command depends upon status of current working directory so that it can add the changes accordingly in the staging area
- Let's see how we can handle this deadlock
- If we observe the initial steps, initially there are only untracked files (after running init for the first time), so, for now we can build our status command so that it will only return untracked files and we will provide other data once the cycle is up and running
- 