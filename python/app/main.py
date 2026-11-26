from __future__ import annotations

import json
import logging
import platform
import secrets
import sys
import time
from collections.abc import Callable
from datetime import UTC, datetime
from typing import Any

from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from app.config import Settings
from app.db import check_postgres

DBChecker = Callable[[str], tuple[bool, str | None]]


class JsonFormatter(logging.Formatter):
    def format(self, record: logging.LogRecord) -> str:
        payload: dict[str, Any] = {
            "time": datetime.now(UTC).isoformat(),
            "level": record.levelname,
            "msg": record.getMessage(),
        }
        for key in ("request_id", "method", "path", "status", "duration_ms"):
            val = getattr(record, key, None)
            if val is not None:
                payload[key] = val
        return json.dumps(payload)


def configure_logging(level_name: str) -> logging.Logger:
    logger = logging.getLogger("tiktokup.python")
    level = getattr(logging, level_name.upper(), logging.INFO)
    logger.setLevel(level)
    if not logger.handlers:
        handler = logging.StreamHandler(sys.stdout)
        handler.setFormatter(JsonFormatter())
        logger.addHandler(handler)
    logger.propagate = False
    return logger


def _envelope(
    *,
    success: bool,
    request_id: str,
    data: Any = None,
    error: dict[str, str] | None = None,
) -> dict[str, Any]:
    body: dict[str, Any] = {
        "success": success,
        "meta": {
            "request_id": request_id,
            "timestamp": datetime.now(UTC).strftime("%Y-%m-%dT%H:%M:%SZ"),
        },
    }
    if data is not None:
        body["data"] = data
    if error is not None:
        body["error"] = error
    return body


def create_app(
    settings: Settings | None = None,
    db_checker: DBChecker | None = None,
) -> FastAPI:
    cfg = settings or Settings.from_env()
    checker = db_checker or check_postgres
    logger = configure_logging(cfg.log_level)
    started_at = datetime.now(UTC)
    started_mono = time.monotonic()

    application = FastAPI(
        title="TikTokUp Python Service",
        version=cfg.app_version,
        docs_url="/docs" if cfg.app_env != "production" else None,
        redoc_url=None,
    )

    @application.middleware("http")
    async def request_context_middleware(request: Request, call_next: Any) -> Any:
        raw_id = request.headers.get("X-Request-ID", "").strip()
        req_id = raw_id if (raw_id and len(raw_id) <= 128) else secrets.token_hex(16)
        request.state.request_id = req_id

        start = time.perf_counter()
        response = await call_next(request)
        duration_ms = round((time.perf_counter() - start) * 1000, 3)

        response.headers["X-Request-ID"] = req_id
        log_level = logging.DEBUG if request.url.path == "/health" else logging.INFO
        logger.log(
            log_level,
            "http_request",
            extra={
                "request_id": req_id,
                "method": request.method,
                "path": request.url.path,
                "status": response.status_code,
                "duration_ms": duration_ms,
            },
        )
        return response

    @application.get("/health")
    def health(request: Request) -> JSONResponse:
        return JSONResponse(
            status_code=200,
            content=_envelope(
                success=True,
                request_id=request.state.request_id,
                data={
                    "status": "ok",
                    "service": cfg.app_name,
                    "environment": cfg.app_env,
                    "uptime_seconds": int(time.monotonic() - started_mono),
                },
            ),
        )

    @application.get("/ready")
    def ready(request: Request) -> JSONResponse:
        db_ok, _ = checker(cfg.database_url)
        checks = {
            "http_server": "ok",
            "postgres": "ok" if db_ok else "error",
        }
        if not db_ok:
            return JSONResponse(
                status_code=503,
                content=_envelope(
                    success=False,
                    request_id=request.state.request_id,
                    data={"status": "not_ready", "checks": checks},
                    error={
                        "code": "SERVICE_UNAVAILABLE",
                        "message": "database connectivity check failed",
                    },
                ),
            )
        return JSONResponse(
            status_code=200,
            content=_envelope(
                success=True,
                request_id=request.state.request_id,
                data={"status": "ready", "checks": checks},
            ),
        )

    @application.get("/info")
    def info(request: Request) -> JSONResponse:
        db_ok, db_version = checker(cfg.database_url)
        db_payload: dict[str, Any] = {
            "engine": "postgresql",
            "connected": db_ok,
        }
        if db_version:
            db_payload["version"] = db_version

        return JSONResponse(
            status_code=200,
            content=_envelope(
                success=True,
                request_id=request.state.request_id,
                data={
                    "service": cfg.app_name,
                    "version": cfg.app_version,
                    "environment": cfg.app_env,
                    "python_version": platform.python_version(),
                    "implementation": platform.python_implementation(),
                    "os": platform.system().lower(),
                    "arch": platform.machine(),
                    "started_at": started_at.strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "uptime_seconds": int(time.monotonic() - started_mono),
                    "database": db_payload,
                },
            ),
        )

    return application


app = create_app()
