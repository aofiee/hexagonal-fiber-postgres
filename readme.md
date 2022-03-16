insert into customers set customers.customer_id = 2,customers.name = 'Arnon',customers.status = 1,customers.city = 'bangkok',customers.zip_code = '10250',customers.created_at = CURRENT_TIME,customers.updated_at = CURRENT_TIME;


insert into customers (created_at,updated_at,customer_id,name,date_of_birth,city,zip_code,status) VALUES (Current_timestamp, Current_timestamp, 2, 'Arnon',Current_timestamp , 'bangkok','10250',1 )