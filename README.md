## Congruit

### The configuration management tool that loves Bash
#### Simple, and lightweight and fully customizable by you

### Table of Contents

**[Description](#Description)**

**[Concepts](#Concepts)**

**[Stockroom](#Stockroom)**

**[Place](#Place)**

**[Work](#Work)**

**[Workplace](#Workplace)**

**[Prerequisites](#Prerequisites)**

**[Usage](#Usage)**

## Description
Congruit is a configuration a lightweight configuration management and automation tool. It is written in Go and works through Bash. It manages shell scripts you created for configure your Linux platforms.

## Concepts
The main concepts of Congruit are

* Stockroom repository
* Works
* Places
* Worksplaces

## Stockroom
The Stockroom are the main repository that describes your platform. Congruit reads the stockroom and does things

## Places
A place is a shell script that must return 0. This means that you are in the right place to do a work.
Example:

Is this Linux server a Centos 7?

```
[ ! -e /etc/redhat-release ] && exit 1
cat /etc/redhat-release | grep "Centos Linux release 7.*"
```
It is nothing of new. They are two lines of shell script

## Works
Works are sheel script that install and configre programs or run Docker containers like in the following example:

```
docker run --rm -p 8888:8080 tomcat:latest &> /dev/null &
```

## Workplace
Workplace are the union between works and places and are Json file.

Example:

```
[
  {
   "places": ["debian","screen_is_not_installed"],
   "works": ["screen_package_apt"]
  },
  {
   "places": ["centos7","screen_is_not_installed"],
   "works": ["screen_package_yum"]
  }
]
```
in the following example, the workplace will take care to decide which is the best strategy for install Screen.
Congruit will execute places and, if they will return 0, do works.

## Prerequisites
1. GO
2. a stockroom. You can take as example the stockroom present in this repo. Please create symlink from stockroom/workplaces/foo to stockroom/workplaces_enabled/foo if you want apply the workplace "foo" during congruit execution

## Usage
1. `git clone https://github.com/lucky-sideburn/congruit.git`
2. `go build conguit.go`

3. Start Congruit

`./congruit  -stockroom-dir=./stockroom`

`./congruit  -stockroom-dir=./stockroom -debug`

