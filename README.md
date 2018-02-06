# graphql-go
Implementation of example project using GraphQL and Golang

## GraphQL schema

This is how the schema is defined:

```
type Temp {
    value: String
}

type Time {
    time: String
}

type Query {
    getCurrentTemperature(city: String): Temp
    getCurrentTime(continent: String, city: String): Time
}
```

There is also a GraphiQL installed. To access it go to http://localhost:8080/graphql
