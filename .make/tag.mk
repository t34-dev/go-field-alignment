################################################################## git-version
GIT_LAST_TAG = @$(shell git tag --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' | head -n 1 || echo "v0.0.0")
GIT_NEW_TAG = $(shell echo $(GIT_LAST_TAG) | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g' | sed 's/^@v//')
CURRENT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

tag:
	@git fetch --tags
	@echo $(GIT_LAST_TAG) | sed 's/^@v//'

tag-new:
	@echo $(GIT_NEW_TAG)

tag-up:
	@if [ "$(CURRENT_BRANCH)" != "main" ]; then \
		echo "Error: You can only create new tags from the 'main' branch."; \
		echo "Current branch: $(CURRENT_BRANCH)"; \
		exit 1; \
	fi
	@git fetch --tags
	@echo $(GIT_NEW_TAG) && git tag "v$(GIT_NEW_TAG)" && git push origin "v$(GIT_NEW_TAG)"
