# btc-api

The api is divided into three directories. The /cmd directory contains the main.go file which starts the program. 
The /configs directory contains the config file. And in /pkg is the main logic of the project. 

The /pkg directory is divided into handler, service and repository. 
The handler is for request handlers, the service is for business logic and the repository is for file management. 

You need an email and password to register. 
The data is stored in a json file, the password is hashed. 
After authorization you get the jwt token, to access the functions, for which you need authorization.
