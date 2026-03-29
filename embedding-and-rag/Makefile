.PHONY: avenv dvenv run list

avenv:
	uv sync

run: avenv
	@if [ -z "$(snippet)" ]; then \
		echo "Usage: make run snippet=<name>"; \
		echo "Run 'make list' to see available snippets."; \
		exit 1; \
	fi
	uv run python3 snippets/$(snippet).py

list:
	@echo "Available snippets:"; \
	for f in snippets/*.py; do echo "  $$(basename $$f .py)"; done

dvenv:
	@echo "uv manages the venv automatically — nothing to deactivate."
