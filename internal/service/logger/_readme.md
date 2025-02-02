**Purpose:** 
DBLogHandler is a custom log handler that writes structured logs into a database. It logs the message, attributes, source (file and line), and timestamp.

**Flow:**
1.  When a log is handled, it extracts the log's message and attributes.
2.    It uses runtime.Caller() to capture the file and line number where the log was triggered.
3.    It marshals the attributes into a JSON string and inserts the log into the database, along with the log level, message, and source.
4.    It only processes logs that meet or exceed the specified minLogLevel.

**Example usage in handler**

Handlers use the `LogHTTPRequestError` in `github.com/brightside-dev/go-chi-rest-api-boilerplate-v2/internal/util` to log errors in the DB.

Handler DB logging flow:
1. Call `LogHTTPRequestError`
2. Get http additional information such as path, ip address and user agent by calling `GetHTTPRequestContext`
3. Map each data into slog attributes using `MapToSlogAttrs`
4. Finally call the slog func `logger.Error(msg, anyAttrs...)` 