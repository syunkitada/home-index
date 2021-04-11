all:
	go run makeindex.go
	git add .
	git commit -am 'update'
	git push
