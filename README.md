### SRE Test
This is a repo with a test task I did one day.
### Description
We ask you to write an application that will display some public information about GitHub users.

### Instructions
- Create a branch for your implementation.
- Design and code simple http server that will provide some public information about requested GitHub user.

Request: GET /github/{username}/repositories

Response: 200 OK

Response example:
```json
{
  "username": "user",
  "public_repositories": 20,
  "followers": 20
}
```
- Expose basic HTTP prometheus metrics (requests count, requests duration) on `/metrics` endpoint and specified port.
- Write the automatic tests for your application. (Consider using Mocks for GitHub API calls)
- Create a Dockerfile that will build and run your application.
- Make a pull-request into the main branch.

### General Requirements
- Golang required.
- Automatic tests required.
- Your code should be formatted and pass included linters (look on the Makefile).
- Port for HTTP server and for metrics endpoint should be configurable through environment variables.
- Usage of external libraries allowed.
- Follow the best practices for writing Dockerfile.
- Keep the commits small.

### Some tips
- We don't expect you to dedicate more than 2-4 hours to this test.
- We don't expect that you will write everything from scratch. You can use API clients, web servers, etc.
- Feel free to choose design patterns you are comfortable with and treat this as a production code.
- You are welcome to ask any questions.
