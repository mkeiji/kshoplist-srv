adminer:
	docker run --name adminer -it -d --network host adminer

mysql:
	docker run --name testdb -it -d --network host\
        -e MYSQL_ROOT_PASSWORD=secret \
        -e MYSQL_DATABASE=testdb \
        mysql

