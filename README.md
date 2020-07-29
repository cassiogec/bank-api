# Go Bank API

This is a small and limited Bank API developed using Go. This API was created with one of the purposes being the study of the language, and for that reason, some of the code in here may not be organized correctly and it may not be in the correct Go standard.

## Run
The code was developed using Go and PostgreSQL, and both of them are running inside docker containers that can be initialized using `docker-compose` with the following code:<br>
* Production: `docker-compose up`
* Tests: `docker-compose -f docker-compose.test.yml up`<br>

Once you have the dockers running, you can call the API through the path `http://127.0.0.1:8080`

## Routes
### `/`
#### `GET /`
Hello Message <br>
Returned Data Example:
```json
{
   "message": "Welcome to the Bank API"
}
```

### `/accounts`

This resource has the following attributes:
* `id`
* `name` 
* `cpf`
* `secret` - It will always be empty because the hash created will be cleared before the return
* `balance` - It has the default value of 10
* `created_at` 

The following routes are accepted:

#### `GET /accounts`
It will return an array with all the register accounts. <br>
Returned Data Example:
```json
[
   {
      "id":1,
      "name":"David Tennant",
      "cpf":"11111111111",
      "secret":"",
      "balance":10,
      "created_at":"2020-01-01T00:00"
   }
]
```

#### `GET /accounts/{account_id}/balance`
It will return a JSON with the balance of the account passed as a parameter. <br>
Returned Data Example:
```json
{
   "account_balance":10
}
```

#### `POST /accounts`
It will create a new account. <br>
Required Body Data Example:
```json
{
   "name":"David Tennant",
   "cpf":"11111111111",
   "secret":"doctor10"
}
```

Returned Data Example:
```json
{
   "id":1,
   "name":"David Tennant",
   "cpf":"11111111111",
   "secret":"",
   "balance":10,
   "created_at":"2020-01-01T00:00"
}
```

### `/login`
Once you have created an account you will need to login to execute some other operations.<br>
The following routes are accepted:

#### `POST /login`
It will return a token and you will have to use it as an authorization.<br>
Required Body Data Example:
```json
{
   "cpf":"11111111111",
   "secret":"secret"
}
```
Returned Data Example:
```json
{
  "token": "YOUR_TOKEN"
}
```

### `/transfers`
In this resource you need to send the token returned in the login resource. It will have to be sent in the header of the HTTP Request in the Authorization field or in the URL of the request.<br>
Example:
* HTTP Header: `Authorization=token: YOUR_TOKEN`
* URL: `http://127.0.0.1:8080/transfers?token=YOUR_TOKEN`<br>

PS: It isn't necessary to send the Account ID for the request paths because it will be used the one hashed in the authorization token. <br>

This resource has the following attributes:
* `id`
* `account_origin_id` 
* `account_destination_id`
* `amount`
* `created_at`

#### `GET /transfers`
It will return an array with all the register tranfers. <br>
Returned Data Example: 
```json
[
   {
      "id":1,
      "account_origin_id":1,
      "account_origin":{
         "id":1,
         "name":"David Tennant",
         "cpf":"11111111111",
         "secret":"",
         "balance":10,
         "created_at":"2020-01-01T00:00"
      },
      "account_destination_id":2,
      "account_destination":{
         "id":2,
         "name":"Tom Baker",
         "cpf":"22222222222",
         "secret":"",
         "balance":10,
         "created_at":"2020-01-01T00:00"
      },
      "amount":5,
      "created_at":"2020-01-02T00:00"
   }
]
```

#### `POST /transfers`
It will return an array with all the register transfers. <br>
Required Body Data Example:
```json
{
   "account_destination_id":2,
   "amount":7.5
}
```

Returned Data Example:
```json
{
   "id":1,
   "account_origin_id":1,
   "account_origin":{
      "id":1,
      "name":"David Tennant",
      "cpf":"11111111111",
      "secret":"",
      "balance":10,
      "created_at":"2020-01-01T00:00"
   },
   "account_destination_id":2,
   "account_destination":{
      "id":2,
      "name":"Tom Baker",
      "cpf":"22222222222",
      "secret":"",
      "balance":10,
      "created_at":"2020-01-01T00:00"
   },
   "amount":5,
   "created_at":"2020-01-02T00:00"
}
```
