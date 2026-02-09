# Go Hash Server

Async file hashing server with worker pool processing.

## Prerequisites

- Go 1.25+

## Run

```bash
go run main.go
```

Server starts on `http://localhost:8080`.

## Test

```bash
go test ./server/handler/ -v
```

## API

| Method | Path | Description |
|--------|------|-------------|
| POST | /upload | Submit a file path for async MD5 hashing |
| POST | /upload_log | Upload a file via multipart form; hashes and stores synchronously |
| GET | /status/{id} | Check job status |
| GET | /hash-content/{id} | Get MD5 hash result |

## Usage

### Async hash (file already on disk)

```bash
mkdir -p /tmp/hash-server/uploads
echo "hello world" > /tmp/hash-server/uploads/test.txt

curl -X POST http://localhost:8080/upload \
  -H "Content-Type: application/json" \
  -d '{"filepath":"/tmp/hash-server/uploads/test.txt"}'
```

Returns a job with status `PENDING`. Poll `/status/{jobId}` until `COMPLETED`, then fetch the hash from `/hash-content/{jobId}`.

### Upload and hash (multipart)

```bash
curl -X POST http://localhost:8080/upload_log -F "file=@/tmp/test.txt"
```

Returns a job with status `COMPLETED` and the MD5 hash immediately.

## Notes

- File paths must be within `/tmp/hash-server/uploads`
- Worker pool size defaults to number of CPU cores
- Jobs that fail hashing are marked as `FAILED`
