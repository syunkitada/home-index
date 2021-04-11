all:
	go run makeindex.go

deploy:
	git add .
	git commit -am 'update'
	git push
