# justfile for ccnow

# run setup on first invocation
default: setup

# setup the local development environment
setup:
	npm install

# build TypeScript
build:
	npm run build

# run tests
test:
	npm test

# run tests in watch mode
test-watch:
	npm test -- --watch

# typecheck the project
typecheck:
	npx tsc --noEmit

# full preflight: build + typecheck + test
preflight: build typecheck test

# build and run with sample input
dev:
	just build
	echo '{"cwd":"/tmp/test","context_window":{"used_percentage":42,"current_usage":{"input_tokens":38000,"cache_creation_input_tokens":2000,"cache_read_input_tokens":1500}}}' | node dist/cli.js

# build and run with tee'd session data (if available)
dev-live:
	just build
	cat tmp/session.json | node dist/cli.js

# remove build artifacts
clean:
	rm -rf dist

# remove everything including node_modules
clobber: clean
	rm -rf node_modules
