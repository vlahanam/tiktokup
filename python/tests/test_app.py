from __future__ import annotations

from fastapi.testclient import TestClient

from app.config import Settings
from app.main import create_app


def _make_settings() -> Settings:
    return Settings(
        app_name="tiktokup-python",
        app_version="0.1.0",
        app_env="test",
        port=8000,
        log_level="warning",
        database_url="postgresql://tiktokup:secret@postgres:5432/tiktokup?sslmode=disable",
    )


def _ok_db(_: str) -> tuple[bool, str | None]:
    return True, "PostgreSQL 18.6"


def _down_db(_: str) -> tuple[bool, str | None]:
    return False, None


def test_health_and_request_id() -> None:
    client = TestClient(create_app(_make_settings(), db_checker=_ok_db))
    res = client.get("/health", headers={"X-Request-ID": "py-trace-001"})

    assert res.status_code == 200
    assert res.headers["X-Request-ID"] == "py-trace-001"
    body = res.json()
    assert body["success"] is True
    assert body["data"]["status"] == "ok"
    assert body["meta"]["request_id"] == "py-trace-001"


def test_ready_when_postgres_ok() -> None:
    client = TestClient(create_app(_make_settings(), db_checker=_ok_db))
    res = client.get("/ready")

    assert res.status_code == 200
    body = res.json()
    assert body["success"] is True
    assert body["data"]["status"] == "ready"
    assert body["data"]["checks"]["postgres"] == "ok"


def test_ready_returns_503_when_postgres_down() -> None:
    client = TestClient(create_app(_make_settings(), db_checker=_down_db))
    res = client.get("/ready")

    assert res.status_code == 503
    body = res.json()
    assert body["success"] is False
    assert body["data"]["checks"]["postgres"] == "error"
    assert body["error"]["code"] == "SERVICE_UNAVAILABLE"


def test_info_returns_python_and_db_metadata() -> None:
    client = TestClient(create_app(_make_settings(), db_checker=_ok_db))
    res = client.get("/info")

    assert res.status_code == 200
    data = res.json()["data"]
    assert data["service"] == "tiktokup-python"
    assert data["python_version"].startswith("3.13")
    assert data["database"]["connected"] is True
    assert data["database"]["version"] == "PostgreSQL 18.6"
