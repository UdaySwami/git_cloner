# git_cloner
Contents:
1] Radio Buttons to connect to GitHub and then clone selected Repo locally.
2] Working positive scenarios

Does Not Contains
1] Unit Tests
2] Complete error handling of API's for negative test cases


To Run 
1] Clone Locally
2]  Execute below command, inside cloned directory.

docker build --tag git_cloner_image1.0 .
docker run --publish 80:80 --detach --name git_cloner git_cloner_image1.0

3] Go to localhost/127.0.0.1 to run.
