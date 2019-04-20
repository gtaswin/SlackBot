# SlackBot
Slack Bot for chat bot automation framework. And its specially for developed for Devops / Release engineer.

## Features
1) RTM based connection.
2) Tightly coupled with shell script for remote execution.
3) Parallel handling of task so each task has own thread.
4) INI file based and dynamic configuration.
5) Admin ID's for secure execution.
6) Regex based reading.
7) more..

## Getting Started
```
./slackbot --config config.ini
```

## Build from source
Using dep package management tool
```
dep ensure
go build -o slackbot
```

## Configuration
Configuration sample file.
```
[main] #MAIN config
token = your-token  #Get the token from slack app.
debug = false #Set true / false for debug option
command = job #Word to execute the shell script
shell = sh #Type of shell

[admin] #Admin ID's to work
user1 = UFBJC6FB1 #User ID 1
user2 = UFBJC6FB8 #User ID 2


[chat] #Words/ Command to do so..
hi = Hello !!! #Normal word
do it = what to do???  
job your ip = sleep 20 && curl ifconfig.me #Shell command to execute by words
```
