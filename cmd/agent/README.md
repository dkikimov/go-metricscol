# Agent

This directory contains **agent** executable code

### Supported settings
* `-a` (env: `ADDRESS` | json: `address`)  **string** \
Address to listen (default "127.0.0.1:8080")
* `-c` (env: `CONFIG`) **string** \
Path to json config
* `-crypto-key` (env: `CRYPTO_KEY` | json: `crypto_key_file_path`) **string** \
Private crypto key for asymmetric encryption
* `-k` (env: `KEY` | json: `hash_key`) **string** \
Key to encrypt metrics
* `-l` (env: `RATE_LIMIT` | json: `rate_limit`) **int** \
Limit the number of requests to the server (default 1)
* `-p` (env: `POLL_INTERVAL` | json: `poll_interval`) **time** \
Interval to poll metrics
* `-r` (env: `REPORT_INTERVAL` | json: `report_interval`) **time** \
Interval to report metrics
