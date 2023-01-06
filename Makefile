docker_test:
	docker run --name test-postgres -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres