adminer:
	docker run --name adminer -it -d --network host adminer

testdb:
	docker run --name testdb -it -d \
		-p 5432:5432 \
        -e POSTGRES_USER=root \
        -e POSTGRES_PASSWORD=secret \
        -e POSTGRES_DB=testdb \
        postgres

testdb2:
	docker run --name testdb2 -it -d \
		-p 3306:3306 \
        -e MYSQL_ROOT_PASSWORD=secret \
        -e MYSQL_DATABASE=testdb \
        mysql

