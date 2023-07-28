---
title: "How to Create Custom Permission Class in Django"
date: 2023-07-25T15:36:25+03:00
tags: ['tutorial', 'python', 'django', 'permissions', 'security']
toc: true
# draft: true
---

Permissions can be a hustle to deal with when developing an api. Suppose we have a number of api views and endpoints where the access permissions are very similar to one another with only slight variations..  
We could create different permission classes with the slight changes that fit the specific endpoint's needs. That works, but that involves a lot of repetition and duplicated code that can be hard to update later down the line.  

This is where an abstract or generic permission class comes in handy. We can define template permissions and logic that can be inherited and extended by child permission classes or the modified in the ViewSet class like so:

```python
# views.py

from somewhere import AbstractPermission

class MyViewSet(viewsets.ModelViewSet):
    serializer_class = MySerializer
    permission_classes = (AbstractPermission,)
    deny_action = {
        "CLIENT": ("create",)
    }
```

### Defining our structure

To get started, we can have an example project for CozyCoin Bank with branches in different countries. Each branch has one manager, several clerks, and marketers. CozyCoin customers are rich and can also have multiple bank accounts. Below are our simple models. I have only included fields necessary for explanations.
All the user models inherit from a default user model class were the common fields for all users are defined. For our specific case, we are using the `roles` field to specify which group the user falls under, but you can use django groups if you wish. I find assigning roles much easier to manage.

```python
# models.py
...
...
class User(AbstractBaseUser):
    class Roles(models.TextChoices):
        MANAGER = "MANAGER", "Manager"
        SUPPORT = "SUPPORT", "Support"
        CLERK = "CLERK", "Clerk"
        CUSTOMER = "CUSTOMER", "Customer"

class BankBranch(models.Model):
    manager = models.ForeignKey(Manager, on_delete=models.CASCADE)

class Customer(User):
    role = models.CharField(max_length=20, default=User.Roles.CUSTOMER.  

class Manager(User):
    role = models.CharField(max_length=20, default=User.Roles.MANAGER.  

class Clerk(User):
    role = models.CharField(max_length=20, default=User.Roles.CLERK.  

class Support(User):
    role = models.CharField(max_length=20, default=User.Roles.SUPPORT.  

class Account(models.Model):
    customer = models.ForeignKey(Customer, on_delete=models.CASCADE)
    bank_branch = models.ForeignKey(Bank, on_delete=models.CASCADE)

class Transactions(models.Model):
    account = models.ForeignKey(Account, on_delete=models.CASCADE)

...
...
```

Let's assume that the serializers and views have also been setup with a sample viewset looking like this:

```python
# views.py
class SomeModelViewSet(viewsets.ModelViewSet):
    serializer_class = SomeModelSerializer
    permission_classes = (SomePermissionClass,) # this is the important part
    def get_queryset(self):
        # your queryset logic here
        pass
```

## Creating the abstract permission class

### A basic permission class

In order to determine if a user has permissions depending on their role, we first have to create an abstract permission class that inherits from rest framework's `BasePermission`.
We will deny any anonymous request since we require the user to be logged in to view their role. If you don't have an account of any sort with CozyCoin, then you shouldn't have access to anything protected by these permissions.

```python
# permissions.py
from rest_framework import permissions

class IsRole(permissions.BasePermission):
    def has_permission(self, request, view):
        if request.user.is_anonymous:
            return Fals.  
    def has_object_permission(self, request, view, obj):
        pass

```

For our permissions to be effective, we first have to determine if the user making the request is in the desired group.
The `role` property is initialized to `None` and is supposed to be modified in the child classes. If it is not, then `SomeError`is raised.
`SomeError` being an error of your choosing. This will ensure that the permissions don't silently fail and deny all requests since the `role` will always be `None` and therefor not match any roles.

```python
# permissions.py
class IsRole(permissions.BasePermission):
    role = None
    def has_permission(self, request, view):
        if request.user.is_anonymous:
            return Fals.  
        return self.has_role(request):

    def has_object_permission(self, request, view, obj):
        return self.has_role(request)

    def has_role(self, request)-> bool:
        if self.role == None:
            raise SomeError
        return request.user.role == self.role

```

And that's it! We can create child classes that specify the roles we want to grant access to. For example:

```pytho.  
# views.py
from some_location.permissions import IsRole

class IsCustomer(IsRole):
    role = "CUSTOMER"

class CustomerOnlyView(viewsets.ViewSet):
    permission_classes = (IsCustomer,)
    ...

```

### Granting only the owner of an object permissions

The permissions we have defined above are perfectly usable, but if you are keen you might have noticed a problem. Taking the `IsCustomer` permissions for example, if the request user is not a customer they will be blocked just fine. The problem is the customer can  request for any customer data that is not theirs and it will be granted to them. This means that any customer can view any other customer's data in a view protected by these permissions. To help with this issue, we need to tweak our `has_object_permission` method.

```python
# permissions.py
    ...
    ...
    def has_object_permission(self, request, view, obj):
        owner = getattr(obj, self.role.lower())
        return owner.id == request.user.id

```

With this small tweak, instead of just giving any data to anyone so long as they have a certain role, we will be able to only limit th.  
data to what they own or is linked to them. This assumes that the foreign key field pointing to the user table is named the same as the role but in all lowercase. This means that if the role is `CLIENT`, then the foreign key field in the model is called `client`. An alternative it to get the name of the fields to check at runtime instead of relying on the role name which could be different.

```python
# permissions
class IsRole(permissions.BasePermission):
    owner_field = None
    ...
    def has_object_permission(self, request, view, obj):
        user = getattr(obj, self.owner_field or self.role.lower()))
        return user.id == request.user.id
  .  
class IsCustomer(IsRole):
    role = "CUSTOMER"
    owner_field = "client"

```

### Granting privileged users access

Wait, what if the person asking for the data has higher access rights like an admin or us? We don't want to make the managers feel less powerful. We have different ways we can deal with this problem. How you approach it will depend on what you named your abstract class and how abstract you want to go. The first approach is to use a variable that grants a set of user roles access. In this example, we can name the field `allow_staff`.  

```python
class IsRole(permissions.BasePermission):
    allow_staff = False
    staff_roles = ["ADMIN", "MANAGER"]
    ...
    def has_permission(self, request, view, obj):
        ...
        if self.allow_staff and request.user.role in staff_roles:
            return True
        ...
    def has_object_permission(self, request, view, obj):
        if self.allow_staff and request.user.role in staff_roles:
            return True
        ...

```

With this approach, we can define the `staff_roles` and set `allow_staff` to true to give the people with `staff_roles` access. If we want to go even more abstract then we can use the second approach were we can define `allow_roles` instead of `staff_roles`. With this we don't need to check if `allow_staff` is true or not, we just check if the current user role is one of the special ones.

### Grant access by view action

Now, that we can filter permissions by roles, we can also try filtering by the action being performed. For example we might not want a support staff member should not be able to view all the customers and their data in one place at once. (This is only useful with a non-incremental system object ID. There's nothing stopping Joe from going up the number of client IDs from one in an incremental system). We might want them to only have access to the specific client asking for support and only for a short period of
time or after they have been per-approved by a client issue ticket handling system. Also when Bob from marketing wants to send an email to all customers over 35, who look like they could use a loan, he will need not any permissions to modify customer data. So we will give him and his team read only access.

To achieve this, we can add it another attribute `deny_actions` to the abstract class. This attribute can be a list containing actions we want to deny on that endpoint. We can then check for the list in the `has_permission` method:

```python
# permissions.py
    ...
    if self.deny_actions and view.action not in self.deny_actions:
        return True
```

This solution is simple, but it will also deny any special privileged users from performing the actions. We can change it up slightly by making `deny_actions` accept a dictionary with the user roles as the keys and denied roles as the values. e.g:

```python
deny_actions = {
    "MARKETING": ["create", "update", "delete"]
}
```

We can put the logic for that in a separate class method and only call it when we need to. We also raise an exception in case the values are not formatted correctly.

```python
# permissions.py
    ...
    ...
    def is_action_denied(self, request, view):
        if isinstance(self.deny_actions, dict):
            actions = self.deny_actions.get(request.user.role)
        elif isinstance(self.deny_actions, list | tuple):
            actions = self.deny_actions
        else:
            raise TypeError(
                f"Invalid type {type(self.deny_actions)} for deny_actions. Expected iterable"
            )
        if actions and view.action in actions:
            return False
        return True
```

### Getting values from the view

Sometimes you do not want to create an entire class just so you can change 1 value and use it on one view only. What would be point of doing all this work? We just need to have a bit more code and logic.  

```python
# permissions.py
    deny_actions = None
    allow_roles = None
    allow_staff = None
    def has_permission(self, request, view):
        deny_actions = getattr(view, 'deny_actions', None)
        allow_roles = getattr(view, 'allow_roles', None)
        allow_staff = getattr(view, 'allow_staff', None)
```

This will allow us to specify the extra values when creating the child permission classes or within the view like so:

```python
# views.py
class SomeModelViewSet(viewsets.ModelViewSet):
    serializer_class = SomeModelSerializer
    permission_classes = (IsClient,)
    allow_staff = True
    deny_actions = {
        "MANAGER": ["create", "update", "delete"]
    }
    def get_queryset(self):
        pass
```

The part where we fetch the filters can also bring up a problem if you have many of them, or if we increase the filters in the future.
So we define a method that is called by `has_permissions` method before doing anything.

```python
# permissions.py
    ...
    def get_or_set_attributes(self, request, view):
    for attribute in dir(self):
        if not attribute.startswith("_"): # ignore any magic methods and private variables
            value = getattr(self, attribute)
            if value == None:
                setattr(self, attribute, getattr(view, attribute, None))
```

## Limitations

So this method of creating things an abstract permissions class is very useful. It means we only have to make a few changes to existing and that's it.
If we need new very specific permissions, we don't have to duplicate some code from somewhere and only change a few things to fit the view. There are some issues with this approach.  

- Too much code in one place.  
As we are trying to make the abstract class as versatile and flexible as possible, we have to cram so much logic and checks in it. This can very quickly
become a nightmare to maintain or to read back in 6 months.  

- Long setup time.  
It took me quite a while to come up with the class and logic. It works well for the most part, but it might be hard to adapt it to a different application.
Is it worth the time it took to make it?  ¯\_(ツ)_/¯ Maybe it will be in the long run as I would just copy paste the class into other projects, and instantly
get customizable permissions without defining them once again.

- Reliance on `role`.  
  If you were to remove the role attribute, major parts of the class will have to be redone or rethought. With this structure, a user can only belong to one role.
  This limitation could be overcome by allowing the primary search field to be set in child class or in the viewset.

- Limited extensibility.  
  There is only so much you can add before it becomes a too much. The best thing at that point would be to create a separate abstract permission class.
- This only works class based views.
  I have not created or tested a function based view version of this code, so it might not work for that use case

## Conclusion

Thanks for getting to the end. Whether you got here by reading or scrolling or copy pasting. You can find the full permission class on [github gits](https://gist.github.com/keystroke3/3166c27ed58abae849f9967d9044db75)
