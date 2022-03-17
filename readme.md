```mysql
insert into customers set customers.customer_id = 2,customers.name = 'Arnon',customers.status = 1,customers.city = 'bangkok',customers.zip_code = '10250',customers.created_at = CURRENT_TIME,customers.updated_at = CURRENT_TIME;
```
```postgres
insert into customers (created_at,updated_at,customer_id,name,date_of_birth,city,zip_code,status) VALUES (Current_timestamp, Current_timestamp, 2, 'Arnon',Current_timestamp , 'bangkok','10250',1 )
```

http://localhost:3000/graph

```graphql
{
  GetCustomer(id: 2) {
    city
    customer_id
    date_of_birth
    name
    status
    zip_code
  }
}


{
  GetCustomers {
    customer_id
    name  
  }
}


{
  Friends: GetCustomers {
    date_of_birth
    name  
  },
  me: GetCustomer(id: 2){
    name
  }
}

mutation {
  CreateCustomer(CustomerID: 3, Name: "Kano", DateOfBirth: "1981-02-14", City: "Bangkok", ZipCode: "10250", Status: 1)
}
```