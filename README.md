# Exercise: Web Scraper for Latest News with REST API

## Objective

The goal of this exercise is to create a web scraper in Go to fetch the latest articles from a news website, 
process the data, and expose it through a REST API. This exercise will help you review core Go concepts like HTTP requests, HTML parsing, 
JSON handling, concurrency, and API design.

## Expected Deliverables

At the end of this exercise, you should have:

1. **A Functional Web Scraper**:

   - Fetches and parses article titles, links, and timestamps from a news website.
   - Uses concurrency to scrape multiple pages in parallel.
   - Handles errors gracefully and logs any failed requests.

2. **A REST API**:

   - Serves the scraped articles in JSON format via a `GET /articles` endpoint.
   - Runs on port `8080`.
   - Returns the following JSON structure:
     ```json
     [
       {
         "title": "Sample Article",
         "link": "https://example.com/article-link",
         "timestamp": "2023-10-15T10:00:00Z"
       },
       ...
     ]
     ```

3. **Enhancement Features**:

   - **Logging Middleware**: Logs each incoming API request with method, path, and timestamp in mongoDB. upload snapshot to S3
   - **Data Refresh**: Automatically re-scrapes the data at regular intervals (e.g., every 5 minutes).

4. **Bonus**:

   - **Rate Limiting**: Implements a simple rate limiter to limit requests to `/articles`.
   - **Pagination**: Supports pagination in the `/articles` endpoint **(if page query string provided)**, e.g., `/articles?page=2`.

## Requirements

### Part 1: Web Scraper

1. - Create an `Article` Table in postgres DB.
2. - Create a `fetchArticles(url string)` function to:
   - Fetch HTML content from a URL.
   - Parse the HTML using `goquery` to extract article information.
   - Append the parsed data to a slice of `Article` structs.
3. Implement concurrency with goroutines to fetch multiple URLs in parallel.
4. Handle errors for HTTP requests and HTML parsing.

### Part 2: REST API

1. Implement a REST API server that:
   - Exposes the scraped data via a `GET /articles` endpoint.
   - Serves JSON data with all scraped articles.
2. Ensure thread-safe access to shared data using a mutex.
3. Verify API functionality with `curl` or Postman.

### Enhancement

- Add middleware for logging.
- Upload logs on S3 weekly
- Add automatic refresh intervals.

### Bonus

- Implement rate limiting.
- Add pagination support for the `/articles` endpoint.

## Submission Checklist

- [✓] **Web Scraper**: Scrapes data from at least two pages and handles errors gracefully.
- [✓] **REST API**: Responds with JSON-formatted article data at `/articles`.
- [✓] **Concurrency**: Uses goroutines to scrape multiple pages concurrently.
- [✓] **Optional Enhancements**: Implements middleware, data refresher
- [✓] **Bonus** (Bonus): rate-limiting, and other optional features as listed above.

## Final Notes

- **Testing**: Test each part thoroughly before moving on to the next section. Use `curl` or Postman for testing API responses.
- **Error Handling**: Ensure that errors are logged clearly and the scraper continues to run even if one URL fails.
- **Documentation**: Include comments explaining key parts of your code, especially for concurrency handling and error handling.
