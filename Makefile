RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

.PHONY: rundev
rundev: ##run web service
	go run app/main.go --target=$(RUN_ARGS)