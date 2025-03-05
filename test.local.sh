go run scripts/loadenv/loadenv.go .env.test.local gotestsum -f testname ${@:-./...}
