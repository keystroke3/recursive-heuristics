---
title: "Sala.pm"
date: 2023-08-14T11:46:54+03:00
tags: ['project']
category: "projects"
toc: true
draft: true
---

Sala is an online scratchpad and clipboard that allows users to quickly note down something or save something from their clipboards.
The clips can then be retrieved later on the same or different devices. Multiple devices can be connected to the same scratchpad at once
allowing for real time data sharing between devices.
Note: **Throughout this document, "scratchpad" and "clipboard" are used interchangeably**

## The Inspiration

My laptop had an issue that caused one of my drives to stop being detected by any linux distribution. Long story short,
I ended up moving to Windows 11 because for some reason that one was detecting the problematic drive.  
With this change, I had another issue. I needed to log into work servers and clone the work repos into Windows Subsystem for Linux (WSL2) but I did not have the RSA pub/priv keys to access these services. There was no easy way for me to create the keys without having access to the servers.
Luckily I had an older laptop where the keys where originally created. All I had to do now was find a way to copy them over to WSL Ubuntu. SSH was not going to work because of how WSL networking worked. The easiest option was to email the files to myself.

Another thing that comes up a lot, is I may be browsing and come across something that I may want to look into on my laptop. In this case, I might make a note in the
notes app and open it up later at my computer. This solves the problem, but I want to do the same from my computer. Save something for when I am on the move.

## The goal

The main goal for Sala was to allow users to save text and files and get to access them on another device without needing to sign in with an account.  
Another goal is to allow users to quickly and easily share what they have saved with others online.  
The service should be simple to use with as little friction as possible.  
The service should be easy to use and respect user privacy.

## Features

With the above goals defined, I came up with a list of features that I would like to have and will help in achieving the goals. Please note that all content marked with work in progress (wip) is either not fully implemented or planned.

### Boards

Scratchpads are where content is saved. Board are divided into private and public boards, public being accessible by anyone with the board's join code and private
requiring a password in addition.

### Realtime multi-directional content saving

All content saved to the board will be appear on all devices actively connected to that board in real time. This allows for easy access to the content.
Any device connected to the board can be used to manage more content, allowing for multi-directional content management.

### End to end encryption (wip)

All content saved to a private board is encrypted using the password provided to the board. The encryption happens on the frontend before it is sent to the backend.  
Decryption also happens on the client so the server has zero knowledge of what content is saved. For technical reasons, the content on public boards cannot be encrypted because of their public nature.

### File sharing (wip)

Since text content allows for a maximum of 5000 characters, it only makes sense to allow users to upload files that hold more content than that. This also allows users
to save binary files such as images, executables and Office documents.

## Architecture

As mentioned before, content is saved in boards, at least from a user perspective. From architectural perspective, the application is broken down into 2: the frontend and the backend.

![Sala Architecture diagram](/sala_architecture.png)

- Backend
    The backend is written in Golang and has a web socket API to handle all board communication and REST API to handle creation and joining requests.
    You can read through the full tour in the [backend section](#backend)

- Frontend
    The frontend is what allows users to interact with the service and is written in typescript and built using Vue3 and Vite.

## Backend

The backend is written in Go and implements the Domain Driven Design with heavy utilization of interfaces for inter-domain communication. The main domains are:

- API
  The API domain handles communication between the client and other internal domains

- Storage
  Handles all transactions that involves data storage both persistent and temporary (key-value storage)

- Http Request Handler
  - Handle Board Creation request
  - Handle Board Joining request
  - Create websocket connections and forward to Message handler
- Message Handler
  - Authentication/Authorization
  - Message validation
  - Message Parsing
  - Message Forwarding
- Session Manager
  - Board initialization and management
  - Content Management
- Clipboard
  - Device Registration and Deregistration
  - Board Password and Name Handling
- Storage
  - Handle database transactions
  - Handle Key-Value storage transactions
