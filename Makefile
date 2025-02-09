
# запуск всех тестов
test:
	make test.unit
	make test.integration

# запуск только юнит тестов
test.unit:
	go test --short -v -count=1 ./...

# запуск только интеграционных тестов
test.integration:
	make test.clean
	docker network create test_network || true
	make postgres.test
	make migrator.test

	go test -v -count=1 ./tests/
	make test.clean

# запуск тествой бд
postgres.test:
	make postgres.test_build
	make postgres.test_run

# билд образа с тестовой бд
postgres.test_build:
	docker build \
	--platform=linux/arm64 \
	-t postgres:test \
	-f postgres.Dockerfile .

# запуск контейнера тестовой бд
postgres.test_run:
	docker run \
	--network test_network \
	--name test_postgres \
	--rm \
	-d \
	-p "54321:5432" \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_PASSWORD=admin \
	-e POSTGRES_DB=scam \
	postgres:test


# запуск миграций для тестовой бд
migrator.test:
	make migrator.test_build
	make migrator.test_run

# билд образа с миграциями для тестовой бд
migrator.test_build:
	docker build \
    --build-arg ENV_FILE=test.env \
    --platform=linux/arm64 \
    -t migrator:test \
    -f migrator.Dockerfile .

# запуск контейнера миграций для тестовой бд
migrator.test_run:
	docker run \
    --restart on-failure \
    --network test_network \
    --name test_migrator \
    migrator:test

test.clean:
	docker stop test_postgres || true
	docker stop test_migrator || true
	docker rm test_migrator || true