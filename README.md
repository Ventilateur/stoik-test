# URL shortener

[Link to test subject](https://github.com/stoikio/jobs/blob/main/go-backend-engineer/README.md)

## Disclaimer

This is a *very* simplified version of an actual URL shortener based on the must-have requirements written in a short 
time. I intentionally did not provide any nice-to-have feature for the time being.

The technique I use for URL shortening is base-62 encoding of a universal counter starting at 1000000000000. 
It ensures each slug is unique but will create issues with scaling.

A few obvious things that I would add if it was a real use case:

- On nice-to-have:
  - Tests: storage tests, domain tests with mock storage and end-to-end tests.
  - Shortening the same URL twice returns the same generated string: a check on database for existing record with the 
    same user/session (user/session id mechanism to be implemented).
  - Expiration dates on URLs: either adding a `expired_on` column and make the query filter on that or use a cronjob to 
    expire the record. Additionally, when using a cache like Redis, also set the expiration on the cache.
  - Clicks counter: a middleware that increment a counter in database before sending the redirection.
- On scaling issue:
  - Cache to speed up look-up process.
  - Database sharding to reduce load. I'd have multiple universal counters with different ranges. A distributed service 
    might be needed to keep track of these ranges.
  - Potential use of a 3-party key generation service.
- Other things:
  - A proper build and database migration process. 
  - CI pipeline.
  - Observability.
  - User/session ID mechanism.

## How to run

From the root of this repo, run:

```shell
make up
```

It will spin up the service along with a postgres database. The service is accessible at `localhost:8080`.

To create a new shortened URL to [https://www.youtube.com/watch?v=dQw4w9WgXcQ](https://www.youtube.com/watch?v=dQw4w9WgXcQ):

```shell
curl --location 'localhost:8080/api/create' \
--header 'Content-Type: application/json' \
--data '{
    "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
}'
```

It will return something like this:

```json
{
  "url":"http://localhost:8080/hBxM5Ag"
}
```

Paste the URL in your browser and see the redirection. The server returns a `302 Found` status with a `Location` points 
to the original URL.
