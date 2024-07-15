![Frame 4](https://github.com/UtkarshM-hub/Bit/assets/70505181/cca9fc7f-c9dd-42c0-83e7-0b4a6a4fe9f9)
bit is a version control system inspired by Git. The main inspiration came from ['Pro Git'](https://git-scm.com/book/en/v2), a book by Scott Chacon (one of the co-founders of GitHub). The 'internals' section of the book provided insightful information about the internal workings of Git and has been a deciding factor in the development of the core workings of Bit.

# Commands
  1. init
     
     initialize a directory as bit directory
     
         bit init .
  3. status
     
     check the status of current working directory and staging area

         bit status
  4. add
     
     add files, changes to staging area inorder to make them ready to commit

     to select all the files

         bit add .

     for single files
     
         bit add <filename>
  5. rm
     
     remove files from staging area. currently only --cached flag is supported
     to select all the files
     
         bit rm --cached .
     
     to select single files
     
         bit rm --cached <filename>
  7. commit
     
     commit the files in staging area

         bit commit -m "<COMMIT_MSG>"
  9. branch
      
     create a branch

         bit branch <BRANCH_NAME>

     view branches or view active branch

         bit branch -a
  11. checkout
      
      switch between branches

          bit checkout <BRANCH_NAME>

# Working of commands
Following information provides just a bird's eye view of how commands are working and it is a bit simplified inorder to make it easy to understand.

## Status Command
![directory structure](https://github.com/UtkarshM-hub/Bit/assets/70505181/64a5f727-e2d0-4122-903e-e30e6177974b)
![Status 1](https://github.com/UtkarshM-hub/Bit/assets/70505181/d4d50dbd-e2ec-4753-ab18-6939c525f61d)
![Status 2](https://github.com/UtkarshM-hub/Bit/assets/70505181/f78f1675-8d52-41e2-8e43-ffe7d343bb95)
![Status 3](https://github.com/UtkarshM-hub/Bit/assets/70505181/6e42eb29-5dee-4df2-aec0-f427412c3972)
![Status 4](https://github.com/UtkarshM-hub/Bit/assets/70505181/b9bd2c1c-c7ae-4d0f-93db-4523a9ab2acc)
![Status 5](https://github.com/UtkarshM-hub/Bit/assets/70505181/31f895e4-3162-4743-bad0-aae7d18da814)
![Status 6](https://github.com/UtkarshM-hub/Bit/assets/70505181/f770c57e-be97-40b5-9f7a-62e81717b31c)

## Add Command
![Add 1](https://github.com/UtkarshM-hub/Bit/assets/70505181/000c6a9c-1b7d-494d-b2a5-6bd20038140b)

## Commit command
![Commit 1](https://github.com/UtkarshM-hub/Bit/assets/70505181/6db089de-cb70-4e06-b1d9-5ae4e27a0c90)
![Commit 2](https://github.com/UtkarshM-hub/Bit/assets/70505181/a06ad118-63c0-416b-a07d-f58c95464b14)
![Commit 3](https://github.com/UtkarshM-hub/Bit/assets/70505181/fa1c6dbd-48e5-4df9-bf38-4463d22d61fa)

## Branch command
![Branch 1](https://github.com/UtkarshM-hub/Bit/assets/70505181/798a256c-6b31-46ec-87b7-11342f05dba5)

## Checkout command
![Checkout 1](https://github.com/UtkarshM-hub/Bit/assets/70505181/01d4aafe-73b0-4e9e-bc82-a10ca1645ea4)
![Checkout 2](https://github.com/UtkarshM-hub/Bit/assets/70505181/9aad07a1-c452-42a9-bfad-9ad35b971915)
![Checkout 4](https://github.com/UtkarshM-hub/Bit/assets/70505181/4aca12c4-f9f2-4a6f-9cbe-b909c7fdf24e)
