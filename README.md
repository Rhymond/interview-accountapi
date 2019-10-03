# Technical decisions

1. Run and quickly analysed `docker-compose.yml` to make sure everything runs on my environment and it don't fail.
2. Download postman templates from [account API](http://api-docs.form3.tech/api.html#organisation-accounts), changed environment in Postman to point to my localhost. Manually examined given fake api using postman and made couple of calls to see how it works.
3. Added golang container to `docker-compose.yml` and created a small test, then run `docker-compose up` to make sure containers are not conflicting and tests runs without issues.
4. Built general API client to make HTTP calls to api. 
5. Built AccountService for API client and implemented `Create`, `Fetch`, `List` and `Delete` operations using Account resource.
6. Implemented integration tests (in `/tests` dir) using gherkin syntax and `godog` package to run them.
7. Linted code using `revive` and fixed linting issues.
8. `git push`

# Tips
* To run all tests including integration tests use `tags`:

```bash
go test --cover -v -tags integration
```


# Exercise

## Instructions

This exercise has been designed to be completed in 4-8 hours. The goal of this exercise is to write a client library 
in Go to access our fake [account API](http://api-docs.form3.tech/api.html#organisation-accounts) service. 

### Should
- Client library should be written in Go
- Document your technical decisions
- Implement the `Create`, `Fetch`, `List` and `Delete` operations on the `accounts` resource. Note that filtering of the List operation is not required, but you should support paging
- Focus on writing full-stack tests that cover the full range of expected and unexpected use-cases
 - Tests can be written in Go idomatic style or in BDD style. Form3 engineers tend to favour BDD. Make sure tests are easy to read
 - If you encounter any problems running the fake accountapi we would encourage you to do some debugging first, 
before reaching out for help

#### Docker-compose

 - Add your solution to the provided docker-compose file
 - We should be able to run `docker-compose up` and see your tests run against the provided account API service 

### Please don't
- Use a code generator to write the client library
- Implement an authentication scheme

## How to submit your exercise
- Create a private repository, copy the `docker-compose` from this repository
- Let us know you've completed the exercise using the link provided at the bottom of the email from our recruitment team
- Usernames of the developers reviewing your code will then be provided for you to grant them access to your private repository
- Put your name in the README
