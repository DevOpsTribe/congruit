## Congruit

### The configuration management tool that loves Bash
#### Simple, lightweight and fully customizable by you

### Table of Contents

**[Description](#description)**

**[Concepts](#concepts)**

**[Stockroom](#stockroom)**

**[Place](#place)**

**[Work](#work)**

**[Workplace](#workplace)**

**[Build your workplaces](#build-your-workplaces)**

**[Prerequisites](#prerequisites)**

**[Usage](#usage)**

## Description
Congruit is a lightweight configuration management and automation tool. It is written in Go but works through Bash. It manages shell scripts you created for configure your Linux platforms.

## Concepts
The main concepts of Congruit are

* Stockroom repository
* Works
* Places
* Worksplaces

## Stockroom
The Stockroom is the main repository that describes your platform. Congruit reads the stockroom and does things

## Place
A place is a shell script that must return 0. You should be in a right place to do a work.
Example:

Is this Linux server a Centos 7?

```
[ ! -e /etc/redhat-release ] && exit 1
cat /etc/redhat-release | grep "Centos Linux release 7.*"
```

## Work
Work is a shell script that installs and configures programs or runs Docker containers like in the following example:

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
the workplace will take care to decide which is the correct strategy for install Screen.
Congruit will execute places and, if they will return 0, do works.

## Build you workplace
1. Create or pull you stockroom

2. You need to describe your places. Example:
  * is this server running a specific Linux Distribution?
  * are there particulare configuration files, software installed, environment variables that describe the role / functionality of this server?
  * places are executed before works... You can inject create files with environment variables that works can use
  Put places in stockroom/places/ folder
 
3. Create works. Put the scripts in stockroom/works. Works do all things like install software, get configuration file from a repository ecc.. 
  
4. Crate workplaces in stockroom/workplaces. Try to make them usable in more environments and follow thi example:

Workplaces are array of hashes
```
[
  {
   "places": ["is_linux", "is_frontend"],
   "works": ["install_apache"]
  },
  
  {
   "places": ["is_linux", "is_frontend", "has_additionl_vhost"],
   "works": ["additional_vhost"]
  },
  
  {
   "places": ["is_production"],
   "works": ["do_backup"]
  }

]
```

## Prerequisites
1. GO
2. a stockroom. You can take as example the stockroom present in this repo. Please create symlink from stockroom/workplaces/foo to stockroom/workplaces_enabled/foo if you want apply the workplace "foo" during congruit execution

## Usage
1. `git clone https://github.com/lucky-sideburn/congruit.git`
2. `go build conguit.go`

3. Start Congruit

`./congruit  -stockroom-dir=./stockroom`

`./congruit  -stockroom-dir=./stockroom -debug`

