# Tests

## Create resource

POST `localhost:8080/`
Payload
```
{
    "long_url": "www.google.com",
    "expiration_date": "2024-03-20"
}
```

Save the response `short_url`.

## Get Resource

GET `localhost:8080/<short_url>`

Try the redirect in a web browser.
By this point, the hits should be at least 2. Let's try to pull the stats. 

## Get Resource stats

GET `localhost:8080/<short_url>/stats`

All time, last 7 days, and last 24 hours should say 2. 

## Delete resource

DELETE `localhost:8080/<short_url>`

You should get a response that the resource was deleted. Try deleting it again. The response will say `resource does not exist`.