# identity-server

Basic rest API for user authentication and registration.
User attributes saved during the registration are:

    a. First Name (mandatory)
    b. Last Name (mandatory)
    c. Password (mandatory)
    d. Email (mandatory)
    e. Company (optional)
    f. Post Code (optional)
    g. Accept Terms (mandatory)
    h. Date created

We have three endpoints in this server:
1) Register cURL command:
 ```
curl -X POST \
   http://localhost:8081/register \
   -H 'Content-Type: application/json' \
   -H 'Host: localhost:8081' \
   -H 'cache-control: no-cache' \
   -d '{
 	"first_name": "John",
 	"last_name": "Doe",
 	"email": "john@gmail.com",
 	"terms": true
 }'
```

2) Login cURL command: 

```
curl -X POST \
  http://localhost:8081/login \
  -H 'Accept: */*' \
  -H 'Content-Type: application/json' \
  -H 'Host: localhost:8081' \
  -d '{
	"email": "john@gmail.com",
	"password": "4CcTZGlU5MnkE0y"
}'
```

3) Home cURL command:

```
curl -X GET \
  http://localhost:8081/home \
  -H 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.l9GD4xkDA_WshvTpHjy08x1VPT8ZnJA9gXpH3CBlIOU' \
  -H 'cache-control: no-cache'
```
