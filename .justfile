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

# auto-format and lint-fix
tidy:
	npx prettier --write .
	npx eslint . --fix

# run format, lint, and type checks (no fix)
check:
	npx prettier --check .
	npx eslint .
	npx tsc --noEmit

# full preflight: build + check + test
preflight: build check test

# remove build artifacts
clean:
	rm -rf dist

# remove everything including node_modules
clobber: clean
	rm -rf node_modules
