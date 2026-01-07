set -e

# Load environment variables from .env if exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Database connection string
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-matchaciee_db}
DB_SSLMODE=${DB_SSLMODE:-disable}

DATABASE_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
MIGRATIONS_PATH="internal/database/migrations"

# Check if migrate CLI is installed
if ! command -v migrate &> /dev/null; then
    echo "Error: 'migrate' CLI not found."
    echo "Please install it using: make install-tools"
    echo "Or manually: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Migration commands
case "$1" in
    up)
        echo "Running migrations up..."
        migrate -path ${MIGRATIONS_PATH} -database "${DATABASE_URL}" up
        echo "Migrations completed successfully"
        ;;
    down)
        echo "Running migrations down..."
        migrate -path ${MIGRATIONS_PATH} -database "${DATABASE_URL}" down
        echo "Migrations rolled back successfully"
        ;;
    force)
        if [ -z "$2" ]; then
            echo "Error: Please specify version number"
            echo "Usage: ./scripts/migrate.sh force <version>"
            exit 1
        fi
        echo "Forcing migration to version $2..."
        migrate -path ${MIGRATIONS_PATH} -database "${DATABASE_URL}" force $2
        echo "Migration forced to version $2"
        ;;
    version)
        echo "Current migration version:"
        migrate -path ${MIGRATIONS_PATH} -database "${DATABASE_URL}" version
        ;;
    drop)
        echo "WARNING: This will drop all tables!"
        read -p "Are you sure? (yes/no): " confirm
        if [ "$confirm" = "yes" ]; then
            migrate -path ${MIGRATIONS_PATH} -database "${DATABASE_URL}" drop -f
            echo "All tables dropped"
        else
            echo "Operation cancelled"
        fi
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Error: Please specify migration name"
            echo "Usage: ./scripts/migrate.sh create <migration_name>"
            exit 1
        fi
        echo "Creating new migration: $2"
        migrate create -ext sql -dir ${MIGRATIONS_PATH} -seq $2
        echo "Migration files created"
        ;;
    *)
        echo "Matchaciee Database Migration Tool"
        echo ""
        echo "Usage: ./scripts/migrate.sh [command]"
        echo ""
        echo "Commands:"
        echo "  up              Run all pending migrations"
        echo "  down            Rollback the last migration"
        echo "  force <version> Force set migration version (use with caution)"
        echo "  version         Show current migration version"
        echo "  drop            Drop all tables (requires confirmation)"
        echo "  create <name>   Create new migration files"
        echo ""
        echo "Examples:"
        echo "  ./scripts/migrate.sh up"
        echo "  ./scripts/migrate.sh down"
        echo "  ./scripts/migrate.sh version"
        echo "  ./scripts/migrate.sh create add_new_field"
        ;;
esac
