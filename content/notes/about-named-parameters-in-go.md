---
title: "Today I learned About Named Returns Values in Go"
date: 2023-08-13T12:12:11+03:00
tags: ['learning']
category: "learning"
toc: true
draft: true
---

In Golang, each function that returns a value must have the return types specified in the function's signature.
Each value that a function says it returns must be returned at some point in the function, otherwise the compiler
will complain about it. If the function in question has many return values, and each of those values must be included
in a return statement that causes the function to exit. This can be very tedious and this where named return values come into play.

### Normal return values

To make things easier to discuss, let us have an example.

Suppose we have want to make structure like so:

```golang

type User

```
