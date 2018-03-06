# Backend Challenge
## Summary

The client class takes the latitude and longitude provided in the route and creates a rudimentary
radius around it. I did not implement the correct mathimatical approach of creating a radius around
these points as there was a lot of maths involving radians and what not, so I simple increased
the maximum and minimum values of the coordinates by 0.01^i where i is the index of the loop,
capped at 9. Each loop will query the database to find any value that contains the search term
and is within the coordinates. If the count is greater than 20, query the database again getting
the values from the table. Then for each row, calculate the distance from the item's location to
the input's location and then sort the list by shortest distance.

## Tools used
[Gorrila Mux](github.com/gorilla/mux) for routing
[Golang Sort](https://golang.org/pkg/sort/) for sorting the items

## Running the code
run ```go run main.go``` in the root folder and make a curl to ```http://localhost:4000/search```

## Running Tests
run ```go test``` in the client folder.
