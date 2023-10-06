---
title: "Docker Quick Tips"
date: 2023-10-06T00:36:09+03:00
tags: ['tutorials']
category: "tutorials"
toc: true
draft: true
---

- [ ] use specific official images - for security
- [ ] use specific versions of images - for better reliability and predictability  
- [ ] use lighter images - for less storage usage and quicker image transfer
- [ ] take advantage of caching through layering - order commands from least frequently to most frequently changing
- [ ] use docker ignore file
- [ ] use multi stage builds
- [ ] Create separate user to run docker containers  
     Create group  
    `RUN groupadd -r <group> && useradd -g <user> <group>`
     Give file ownership to the user
    `RUN chown -R <group>:<user> /app`
