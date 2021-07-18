SHELL := /bin/bash
.PHONY: check

.SILENT: build getlasttag

# build: getlasttag
# 	$(info +Компиляция)
# 	go build -o ./bin/main/main cmd/app/main/main.go
build:
	$(info +Компиляция)
	go build -o ./bin/main/canvas cmd/main/main.go

buildwin:
	$(info +Компиляция windows)
	GOOS=windows GOARCH=amd64 go build -o ./bin/main/wincanvas.exe cmd/app/main/main.go

# run: build
# 	$(info +Запуск)
# 	./bin/main/main -qwert -89  copy -r r234 -w ./e help8 help -help -mm  reverse  addPath -p -44.89 -p -788 -p 879 

run: build
	$(info +Запуск)
	./bin/main/canvas

getlasttag:
	git describe --tags
# make gittag tag=vx.x.x
gittag: check

check: 
#ifndef $(tag)#"$(git describe --tags)"; 
	
	@{ \
	set -e ;\
	line=`git describe --tags`;\
#	echo $$line; \
	echo Введите новый tag? последний тег: $$line [n - отмена];\
	read line;\
	if [[ $$line == "n" ]]; \
	then \
	echo вы отказались; \
	exit 7;\
	else \
	git tag $$line ;\
	git push origin --tags ;\
	echo end;\
	fi;\
	}
#endif
	
#	@git tag $$line
#	@git push origin --tags

help:
	go doc -all ./internal
	go doc -all ./cmd/app/main
	go doc -all ./pkg/structs