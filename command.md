docker compose \
  -f docker-compose.yml \
  -f docker-compose.dev.yml \
  up

docker compose -f docker-compose.dev.yml up --build // jika app belum ada


# Change Docker Compose tambahkan Build jika ingin deploy ke railways

docker compose run --entrypoint sh app //masuk entry point

cd /app
ls -la

docker compose \
  -f docker-compose.yml \
  -f docker-compose.dev.yml \
  down -v // hati2 -v menghapus data volumes

docker compose \
  -f docker-compose.yml \
  -f docker-compose.dev.yml \
  build --no-cache

docker compose \
  -f docker-compose.yml \
  -f docker-compose.dev.yml \
  up

go test ./...

go test -cover ./...

Karena ini unit test, bukan integration test.