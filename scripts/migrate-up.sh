echo "Миграция UP (init-up.sql)..."

PGPASSWORD="rootpassword123" psql -h "localhost" -p "5434" -U "root" -d "RIP" -f init-up.sql

echo "Миграция UP выполнена."