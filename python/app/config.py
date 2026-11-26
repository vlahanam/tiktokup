from __future__ import annotations

import os
from dataclasses import dataclass
from urllib.parse import quote, urlencode


@dataclass(frozen=True)
class Settings:
    app_name: str
    app_version: str
    app_env: str
    port: int
    log_level: str
    database_url: str

    @classmethod
    def from_env(cls) -> Settings:
        database_url = os.getenv("DATABASE_URL", "").strip()
        if not database_url:
            host = os.getenv("DB_HOST", "postgres").strip() or "postgres"
            port = os.getenv("DB_PORT", "5432").strip() or "5432"
            user = quote(os.getenv("DB_USER", "tiktokup").strip() or "tiktokup", safe="")
            password = quote(
                os.getenv("DB_PASSWORD", "tiktokup_secret").strip() or "tiktokup_secret",
                safe="",
            )
            name = quote(os.getenv("DB_NAME", "tiktokup").strip() or "tiktokup", safe="")
            sslmode = os.getenv("DB_SSLMODE", "disable").strip() or "disable"
            query = urlencode({"sslmode": sslmode})
            database_url = f"postgresql://{user}:{password}@{host}:{port}/{name}?{query}"

        raw_port = os.getenv("PORT", "8000").strip() or "8000"
        try:
            port_num = int(raw_port)
        except ValueError:
            port_num = 8000

        return cls(
            app_name=os.getenv("PYTHON_APP_NAME", "tiktokup-python").strip() or "tiktokup-python",
            app_version=os.getenv("APP_VERSION", "0.1.0").strip() or "0.1.0",
            app_env=os.getenv("APP_ENV", "development").strip() or "development",
            port=port_num,
            log_level=(os.getenv("LOG_LEVEL", "info").strip() or "info").lower(),
            database_url=database_url,
        )
