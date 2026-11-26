"""Example standalone script runnable via Docker (`make python-run SCRIPT=scripts/example.py`)."""

from __future__ import annotations

import json
import platform

from app.config import Settings
from app.db import check_postgres


def main() -> None:
    cfg = Settings.from_env()
    db_ok, db_version = check_postgres(cfg.database_url)
    print(
        json.dumps(
            {
                "service": cfg.app_name,
                "python_version": platform.python_version(),
                "postgres_connected": db_ok,
                "postgres_version": db_version,
            },
            indent=2,
        )
    )


if __name__ == "__main__":
    main()
