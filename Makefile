# Makefile

# test package
PACKAGE=./app/handlers

# coverage files
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

# target to run tests and collect coverage
coverage:
	@echo "Running tests and collecting coverage..."
	go test -coverpkg=./... -coverprofile=$(COVERAGE_OUT) $(PACKAGE)
	@echo "Generating HTML report..."
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Opening coverage report in Firefox..."
	firefox $(COVERAGE_HTML)

# clean
clean:
	@echo "Cleaning up coverage files..."
	rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)

guid:
	@uuid=$$(uuidgen); \
	echo $$uuid; \
	echo -n $$uuid | xclip -selection clipboard;

.PHONY: coverage clean giud
