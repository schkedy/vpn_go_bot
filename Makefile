.PHONY: run

run:
	go run cmd/server/main.go

git :
	git add .
	git commit -m "$(text)"
	git push origin main