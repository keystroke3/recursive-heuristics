---
title: "A Different Way to Think About Recursion"
date: 2023-11-22T13:21:34+03:00
tags: ['tutorials', 'recursion', 'python', 'golang']
category: "tutorials"
toc: true
draft: true
---

Example:
[Example question](https://www.boot.dev/assignments/5dfdf265-e5df-4f16-a6bf-6aca05ebb9a7)

> Those who don't understand recursion are doomed to repeat it.

Recursion is a simple but very confusing concept. When I first learned about it, I sure was confused as to how a function can call itself. I was working on some
code the other day and I had to use recursion to solve a problem, and while I was writing the function, I had a realization that changed how I look at recursion. I would like to
share that with you here today.

First, let's start by a small explanation of what recursion is incase this is new to you. If you already know what it is, you can skip to the next section.

## What is recursion?

Recursion is a programming concept that arises from the Functional Programming paradigm. The world of functional programming is built on a few key principles:

- Minimal to no mutation of external states. Everything that a function does must be contained within itself and there should be no side effects. This means a function should only rely on it's inputs to generate its outputs and not manipulate any external variables or entities including their inputs absolutely necessary.
- Functions must be deterministic. Given the same inputs, the same outputs should be expected. This is because there are no side effects as mentioned above.
- All functions must return something. Functions have inputs and produces outputs and since there should not be any side effects, the best way to know the output of a function is to read the function's return values.
A function that adheres to the the first 2 principals is known as a *pure function*. Pure functions are stateless in design because they do not rely on or manipulate external states.

Ok so with those core pillars of functional programming in mind, we can tackle recursion. To answer the main question, here is my simple definition of recursion:

**Recursion is when a function calls itself in its with different inputs until a base case is reached**.  
It's not a perfect definition, but it captures everything that is needed to understand this article. A recursive function, if written correctly, is an example of a pure function. It takes some input and iterates over it and produces some output and does not mutate any external states. Recursive functions must also have a return value otherwise the loop will never end, or in the case of Python, you would reach the maximum allowed recursion depth.

### Example

A simple example of a recursive function is a factorial calculator function:

```python
def factorial(n:int) -> int:
    if n == 1:
        return 1
    return n * factorial(n-1)
```

<!-- markdownlint-disable MD037 -->
A factorial is a product of all integers from 1 up to a limit N, `N * (N-1 * (N-2 * .... (N-(N-1))))`.
<!-- markdownlint-enable MD037 -->

The function does exactly this. It takes an integer `n` and multiplies it by the the return value of calling factorial with the number just below `n`. The base case in this case, is when we reach `1` which is `n-(n-1)`. At this point we know we have gone through all the values from where we started.

## The Rethinking

When I was learning recursion, thought that all recursive functions must have the same stature as the factorial function shown above. I.e., a simple base case with an `if` statement, and recursive call being at the end with a `return` statement attached to it. This worked for me most part, with simple easy to understand logic, but it quickly falls apart when applied to anything more complex.
Recursion is most powerful with Nested and tree like data structures, so something had to change with this thiking

I was tasked with traversing a tree-like data structure and getting all paths to all the leaves. This is akin to traversing the file structure of inside a folder and getting all the paths to the files at the end of each directory. Here is an example of such a file structure in the form of a Python dictionary:

```python
# Input
{
    "Documents": {
        "Proposal.docx": None,
        "Report": {"AnnualReport.pdf": None, "Financials.xlsx": None},
    },
    "Downloads": {"picture1.jpg": None, "picture2.jpg": None},
}
# Expected Output:
[
    "/Documents/Proposal.docx",
    "/Documents/Report/AnnualReport.pdf",
    "/Documents/Report/Financials.xlsx",
    "/Downloads/picture1.jpg",
    "/Downloads/picture2.jpg",
]
```

The nodes in this case are keys that have other dicts as their values, and the leaves were those that have `None` as their values.

### First realization

**The recursive call does not have to be at the end or involve a return statement.**

Here is the one of many attempt at solving the problem described above:

```python
def list_files(current_node, current_path="", paths=[]):
    for k,v in current_node.items():
        new_path = f"{current_path}/{k}"
        if v == None:
            paths.append(new_path)
            return paths
        return list_files(v, new_path, paths)

```

A problem with this approach is the is an early `return` when checking the base case of whether we have reached the end of a branch. This results in only the first branch path being returned. So we replace it with `continue`.

This change resulted in `None` being returned because of course, the function itself is not returning anything. When each iteration of the for loop reaches the end of a branch, it will exit and the return statement will not be reached.

It is at this point that I realized, we are supposed to return `paths`, not the results of `list_files()` and it became clear that you don't have to have the recursive call at the end with with a `return`:

```python
def list_files(current_node, current_path="", paths = []):
    for k,v in current_node.items(): 
        new_path = f"{current_path}/{k}"
        if v == None:
            paths.append(new_path)
            continue # replace `return paths`
        else:
            return list_files(v, new_path, paths)
    return paths
```

Finally it worked and I was able to get the expected results, but of course, it had another problem. Pause here to see if you can spot it. *Hint: Look back at the pillars of functional programming*

### Second Realization

**Recursion is about simplifying going forward and clarifying going backwards**  
Let's look at that factorial description again:
$$N \times (N-1 \times (N-2 \times .... (N-(N-1))))$$

Whether you use $PEMDAS$ or $BODMAS$, the same thing always comes first in the order of operations: Parenthesis. Therefor, if we were to evaluate the expression for a factorial, we would have to start from the right going on the left. Basically start at the end and work our way back to the beginning.

Here is an example analogy. Imagine the CEO of a large tech company wants a status report of the whole company for her to report to the investors. The CEO calls a meeting with the other executives of the different arms of the company, lets say the CFO and CTO, and asks for the reports. The CTO calls the departmental managers below him and asks for a similar report, and the managers go the leaders of each team and ask for the same thing. Lastly each team member reports to what they are working on and the status of the projects from their team leader.  
The team members are at the bottom of the chain with the actual information. They don't need to ask anyone else below them and are the "base case". So now each team member gives a report which the team leader compiles together with a report of their own activities and sends it up to the managers. The managers compile all reports from the teams they manage and add information about their office and push it up. This goes on until we finally get back to the CEO who then compiles all this information and gives it to the investors. The investors can then decide how private jets they should get for their spouses.

When we no longer have any more uncertainty about what information the function should return, we have reached the bottom of the tree, then core of the nesting, the base case.
From here, it is about compiling the information going up. Each stage providing more and more certainty to the caller. Now the tricky part is deciding how to 'compile' this information, and what information is uncertain and needs to be clarified by 'lesser' function call.
So going down the recursion stack trace, we are simplifying the input and instructions, and going up, we are providing clarification on uncertainty.

To go back to the factorial, we see that at each stage, we are unsure of what the total product of the numbers that come before are. So we simplify the input by taking away the current value of $N$ until we have certainty that the product of 1 and itself is 1.

```python
def list_files(current_node, current_path=""):
    file_list = []
    for k, v in current_node.items():
        if v != None:
            file_list.extend(list_files(v, f"{current_path}/{k}")) 
        else:
            file_list.append(f"{current_path}/{k}")
    return file_list

```

The

