test:
	go test -v

demo:
	go run example/*

setup-pre-commit:
	brew install pre-commit
	pre-commit install
	
