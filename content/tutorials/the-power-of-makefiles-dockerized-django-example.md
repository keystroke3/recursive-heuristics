---
title: "The Power of Makefiles.A Dockerized Django Example"
date: 2023-10-05T22:30:33+03:00
tags: ['tutorials', 'django', 'makefile', 'docker', 'bash', 'shell']
category: "tutorials"
toc: true
draft: true
---

Makefiles are very useful files when it comes to automation and can help save a lot of time with repeptitive tasks during development. In this article we will have a quick introductory look at makefiles and look at an example using Django and Docker.

## What are makefiles

Makefiles are a part of the Make automation tool that is used to build software, usually binary, from source. Make reads instructions from makefiles to know how to perform the compilation. Make is language agnostic and will work so long as the sytax in the make file is correct and the commands issued work correctly. Make itself does not compile software, it just automates the process by running a series of specified commands. If the commands execute successfully, then the software is compliled.

## how they work
Make files follow a simple sytax for basic usage. It is just `make [target]` where target is the specific instruction set to execute. If no target is given, make will execute the first target in the makefile. Usually, it is a function called
- [ ] explain the django project
- [ ] explain docker
- [ ] construct the makefile line-by-line
- [ ] makefile vs shell script
- [ ] makefile vs just file
