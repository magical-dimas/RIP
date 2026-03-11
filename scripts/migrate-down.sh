echo "Миграция DOWN (init-down.sql)..."

PGPASSWORD="rootpassword123" psql -h "localhost" -p "5434" -U "root" -d "RIP" -f init-down.sql

echo "Миграция DOWN выполнена."