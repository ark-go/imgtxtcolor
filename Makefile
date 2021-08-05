SHELL := /bin/bash
.PHONY: check

.SILENT: build getlasttag buildzip


build:
	$(info +Компиляция)
	go build -o ./bin/main/canvas cmd/main/main.go
buildzip:
	$(info +Компиляция с жатием)
	go build -ldflags "-s -w" -o ./bin/main/canvas cmd/main/main.go
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

gitsend: gitsave check

# make gittag tag=vx.x.x
gittag: check

check: 
#ifndef $(tag)#"$(git describe --tags)"; 
	
	@{ \
	set -e ;\
	line=`git describe --tags | cut -d "-" -f 1`;\
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
	
gitsave:
	@{ \
	set -e ;\
	git status --short;\
	line=`git describe --tags | cut -d "-" -f 1`;\
	echo Отправляем изменения на сервр, последний тэг: $$line; \
	read -p "Введите комментарий [ n-отменить отправку ( по умолчанию: $$line )]: " commitname;\
	line=$${commitname:-$$line};\
	if [[ $$commitname == "n" ]]; \
	then \
	echo вы отказались отправлять изменения; \
	exit 0;\
	else \
	git commit -a -m $$line ;\
	git push origin;\
#	git push origin --tags ;\
	echo end;\
	fi;\
	}

#	check

#	@git tag $$line
#	@git push origin --tags

help:
	$(info run - соберем и запустим)
	$(info build - соберем без запуска)
	$(info gitsend - отправим на сервер, поставим тэг)
	$(info gitsave - отправим на сервер)
	$(info gittag - установим новый тэг)
	$(info gitlasttag - показать последний тэг)

#	go doc -all ./internal
#	go doc -all ./cmd/app/main
#	go doc -all ./pkg/structs