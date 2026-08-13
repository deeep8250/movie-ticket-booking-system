$env:DB_HOST="127.0.0.1"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"
$env:DB_NAME="movie_booking"

$env:REDIS_HOST="127.0.0.1"
$env:REDIS_PORT="6379"
$env:REDIS_PASSWORD="redis_password_12234"

$env:JWT_SECRET="its-secret-bro"

go test ./tests/integration -v -count=1