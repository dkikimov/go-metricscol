# Server

This directory contains **server** executable code

### Supported settings
* `-a` (env: `ADDRESS` | json: `address`)  **string** \
  Address to listen (default "127.0.0.1:8080")
* `-c` (env: `CONFIG`) **string** \
  Path to json config
* `-crypto-key` (env: `CRYPTO_KEY` | json: `crypto_key_file_path`) **string** \
  Private crypto key for asymmetric encryption
* `-d` (env: `DATABASE_DSN` | json: `database_dsn`) **string** \
    Database DSN
* `-f` (env: `STORE_FILE` | json: `store_file`) **string** \
File to store metrics (default "/tmp/devops-metrics-db.json")
* `-i` (env: `STORE_INTERVAL` | json: `store_interval`) **time** \
    Interval to store metrics
* `-k` (env: `KEY` | json: `hash_key`) **string** \
  Key to encrypt metrics
*  `-r` (env: `RESTORE` | json: `restore`) \
Restore metrics from file (default true)
* `-t` (env: `TRUSTED_SUBNET` | json: `trusted_subnet`) **string** \
  Trusted subnet

