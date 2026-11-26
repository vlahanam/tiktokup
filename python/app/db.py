from __future__ import annotations

import psycopg


def check_postgres(dsn: str, timeout_seconds: int = 2) -> tuple[bool, str | None]:
    """Verify PostgreSQL connectivity and return a concise version string."""
    try:
        with psycopg.connect(dsn, connect_timeout=timeout_seconds) as conn:
            with conn.cursor() as cur:
                cur.execute("SELECT version()")
                row = cur.fetchone()
                if not row or not row[0]:
                    return True, None
                tokens = str(row[0]).split()
                version = " ".join(tokens[:2]) if len(tokens) >= 2 else str(row[0])
                return True, version
    except Exception:
        return False, None
