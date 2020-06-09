# Rest API  for Streaming log to google pub/sub
it will be http post data send to pub/sub topic

## Message classification
Using url path
ex)

```bash
curl 127.0.0.1:8080/bigqueryproject/dataset/table -d "{\"col1\":\"data1\",\"col2\":123}"
```

it will be send pub/sub to attribute `map[string]{"Path":"/bigqueryproject/dataset/table"}` and message `[]byte("{\"col1\":\"data1\",\"col2\":123}")`


## After collecting
Cleaning and Save pub/sub message as gcs file using dataflow and load to bigquery
