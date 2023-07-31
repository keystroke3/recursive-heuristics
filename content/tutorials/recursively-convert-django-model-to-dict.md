---
title: "How to Recursively Convert Django Model to Dict"
date: 2023-07-25T16:15:35+03:00
tags: ['python', 'django', 'tutorial', 'json']
category: "tutorial"
toc: true
---

In this article we will look at how to create your own recursive serializer in Django without using the Django Rest Framework
Sometimes when working with Django, we may have some data that we want to  serialize (convert to JSON) but we do not have the option of
using the serializers that come with Django Rest Framework. The data can also take the same form for many cases and writing a new serializer
for all of those cases can be tedious and repetitive. This sounds like something a utility function that can be used to generalize the logic and be reused where needed, hence this article's existence.

By the end of the article, we should have a utility function like this:

```python
import json

from somewhere import ModelFoo
from utilities import generic_serializer

instances = ModelFoo.objects.filter(some_filter)
serialized_instances = generic_serializer(instance, i_fields=['bar'], i_models=['baz'])

with open('Foo_data.json', 'w') as fd:
    json.dump(serialized_instances, fd, default=str, indent=2)

```

## Defining the models

To explain things better, let us use an example. Suppose we have a Django application that stores data for a chain of car dealerships.
Each dealership has cars that customers go to one of them to buy a car. The example `models.py` is shown below:

```python
# models.py
from django.db import models

class Dealership(models.Model):
    name = models.CharField(max_length=20) 
    street_name = models.CharField(max_length=20) 
    zip_code = models.CharField(max_length=6) 

    def __str__(self):
        return self.name

class Car(models.Model):
    name = models.CharField(max_length=20) 
    brand = models.CharField(max_length=20) 
    year = models.CharField(max_length=20) 
    price = models.DecimalField(max_digits=7, decimal_places=2)
    dealership = models.ForeignKey(Dealership, on_delete=models.CASCADE)

    def __str__(self):
        return f"{self.name} {self.brand}"


```

## Creating the generic serializer

### Finding the fields

To start us off the first thing our serializer need to be able to do to create a `dict` of  the `fields` and `values` of our given model instance. In this article, we will focus on the
`Car` model instance and try to serialize it. To give us something to work with we create some sample
instances:

```python
# backups.py
dealership = Dealership.objects.create(name="LopezCars", street_name="Tarmac", zip_code="90210")
# {"id":1, "name":"LopezCars", "street_name": "Tarmac", "zip_code": "90210"}

car = Car.objects.create(name="Camión ", brand="Suave", year="2020", price="2283", dealership=1)
# {"id":1, "name":"Camión", "brand":"Suave", "year":"2020", "price":Decimal('2283.00'), "dealership": 1}
```

The simplest way to convert a model instance to a dict in django is to use the `model_to_dict` function that is built into django.
To use it, simply import it and pass a model instance to it. It supports include all or exclude all but the specified fields.

```python
>>> from django.forms import model_to_dict
>>> from .backups import car
>>>
>>> model_to_dict(car)
>>> {"name":"Camión", "model":"Suave", "year":"2020", "price":Decimal('2283.00'), "dealership": 1}
>>>
>>> model_to_dict(car, fields=["brand", "price"])
>>> {"model":"Suave", "price":Decimal('2283.00')}
>>>
>>> model_to_dict(car, exclude=["brand", "price"])
>>> {"name":"Camión", "year":"2020", "dealership": 1}

```

While `model_to_dict` is very useful for a simple model, it is not powerful enough for what we may need, like recursively fetching related models. Django models have the `_meta` API  which is an instance an `django.db.models.options.Options` object that allows us to fetch all the field instances of a model. One of the properties available in `_meta` is the `fields` which is a `django.utils.datastructures.ImmutableList`.

We can use this list to get all fields and construct a dict object out of it like so:

```python
# backups.py
...
fields = {}
for f in car._meta.fields:
    fields[f.name] = getattr(model, f.name)

# if you wish to use a comprehension:
fields = {f.name: getattr(model, f.name) for f in car._meta.fields}

# {"id":1, "name":"Camión", "brand":"Suave", "year":"2020", "price":Decimal('2283.00'), "dealership": <Dealership: LopezCars> }
```

Looking at the output, you may notice that with this method, we preserve the object Primary Key (id) a.k.a pk and the `dealership` is not just the id. If we do the same thing for `dealership`:

```python
# backups.py
...
# if you wish to use a comprehension:
fields = {f.name: getattr(model, f.name) for f in dealership._meta.fields}

# {"id":1, "name":"LopezCars", "street_name": "Tarmac", "zip_code": "90210"}
```

And now we combine the two:

```python
# backups.py
fields = {f.name: getattr(model, f.name) for f in car._meta.fields}
fields['dealership'] = {f.name: getattr(model, f.name) for f in dealership._meta.fields}

# {"id":1, "name":"Camión", "brand":"Suave", "year":"2020", "price":Decimal('2283.00'), "dealership": {"id":1, "name":"LopezCars", "street_name": "Tarmac", "zip_code": "90210"}}
```

We have successfully converted the model into a dictionary. This is nice, but it is specific to `car` so to make it generic using the magic of recursion:

```python
# utilities.py
from django.db.models import Model
def get_fields(model: models.Model, fields: dict = None) -> dict:
    if not fields:
        fields = {}
        for f in model._meta.fields:
            fields[f.name] = getattr(model, f.name)

    for name, value in fields.items():
        if not isinstance(value, model.Model):
            # skip non-relational (ForeignKey) fields
            continue
        fields[name] = getattr(value, "pk") if name in ignore_models else get_fields(value)
    return fields

get_fields(car)

# {"id":1, "name":"Camión", "brand":"Suave", "year":"2020", "price":Decimal('2283.00'), "dealership": {"id":1, "name":"LopezCars", "street_name": "Tarmac", "zip_code": "90210"}}
```

Now we have a generic function to transform any model into a dict object and that can recurse into the other models within it. If that is all you needed, then you can stop here.
We have made a slightly more powerful version of `model_to_dict`.

### Extra Functionality

The `get_fields` function is fine the way it is and is perfectly usable now. But what if we don't need the to see the `brand` field or maybe `dealership` related model? If we want to get information for more than one car we would have to run the function on both of them. To achieve these things we need to expand it with a few things.

```python
# utilities.py
from django.db.models import Model

def generic_serializer(
    instances: list,
    exclude_models: list = None,
    exclude_fields: list = None,
) -> list:
    def get_fields(model: Model = None, fields: dict = None) -> dict:
        if not fields:
            fields = {}
            for f in model._meta.fields:
                if (
                    not f.is_relation
                    and f.name in exclude_fields
                ):
                    continue
                fields[f.name] = getattr(model, f.name)

        for name, value in fields.items():
            if not isinstance(value, Model):
                continue
            # Replace the excluded related model field with their primary key value
            # Instead of storing the `model.__str__()` value, it shows the model.pk
            # if it is in `exclude_models`
            fields[name] = getattr(value, "pk") if name in exclude_models else get_fields(value)
        return fields

    return [get_fields(instance) for instance in instances]

```

To check if a field `f`  is a relation field i.e `models.ForeignKey` field, we can use `f.is_relation`.  For the `car` model the `dealership`
field is a relational field to `Dealership` models. Therefor, the value for for `dealership.is_relation` is `True`. This allows the exclusion of non-relation field in `exclude_field` without excluding a relation field with the same name and vice versa. The line `if not is instance(value, Model)` checks if there there are any relation fields to trigger a recursion.

We are done! (^O^)/ It is a simple function, but it is very useful.

Thanks for taking the time to read this article and I hope it has been useful to you.
