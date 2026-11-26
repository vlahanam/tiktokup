# syntax=docker/dockerfile:1

FROM python:3.13-slim AS base
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    VENV_PATH=/opt/venv \
    PYTHONPATH=/app
ENV PATH="${VENV_PATH}/bin:${PATH}"
WORKDIR /app

FROM base AS builder
RUN python -m venv "${VENV_PATH}"
COPY pyproject.toml ./
COPY app ./app
RUN pip install --upgrade pip && pip install .

FROM builder AS dev
RUN pip install -e ".[dev]"
COPY . .
EXPOSE 8000
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]

FROM base AS production
RUN groupadd --system appgroup && useradd --system --gid appgroup appuser
COPY --from=builder "${VENV_PATH}" "${VENV_PATH}"
COPY app ./app
COPY scripts ./scripts
USER appuser
EXPOSE 8000
HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://127.0.0.1:8000/health', timeout=2)" || exit 1
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
