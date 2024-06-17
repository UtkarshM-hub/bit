# status command

<b>description:</b>
provides an overview of the current state of the working directory and staging area. It shows which changes have been staged, which haven't, and which files aren't being tracked by Git

## Understanding the command
- it is basically showing us the 

    - files which are untracked
    - files which are modified
    - files which have been deleted
    - files which are tracked and ready to commit

## Think about it
- So it is all about tracking right, if we are able to track the files in working directory and compare them with those stored in staging area then we will be able to implement the above things.
- How can we track the files and on what basis we can decide if it is new, if it is modified, if it is deleted or if it is tracked
- We can use following logic to filter the files accordingly

    - <b>New file</b>: 
        - If the entry for current file doesn't exist in staging area  
    - <b>Modified file</b>: 

        - We will have to check both the last modified time and content of the file inorder to determine if it is modified or not
        - If last modified time of file and time which is stored in staging area is not equal and if the content is different then it is updated otherwise it is not updated
    - <b>Deleted file</b>:

        - If an entry exists for a file but it is not present inside the directory 
    - <b>Tracked file</b>:

        - If files entry exists staging area and it is not modified
- By observing above criteris of identification we require following things to keep track of files

    - modified at
    - content of the file

- Now the question comes that how are we going to keep track of files in our staging area? or how can we represent the staging area?
- For that git uses index file which is at the root of .git folder
- Its content is in binary format but for our simplicity we will keep ours in plain text just to understand the behaviour and working of our application
- Now the thing over here is, creating or updating the entry from index file is done by using add command ( which we haven't created till now ), it basically adds the file in staging area
- Let's impelement the add command
- 