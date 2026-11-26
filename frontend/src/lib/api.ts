import type {
  ApiEnvelope,
  HealthData,
  ReadyData,
  SystemInfoData,
} from "@/types/api";

export interface ApiCallResult<T> {
  envelope: ApiEnvelope<T>;
  status: number;
  latencyMs: number;
  requestId: string | null;
}

function getApiBaseUrl(): string {
  if (typeof window === "undefined") {
    const internal = process.env.INTERNAL_API_URL || "http://backend:8080";
    return `${internal.replace(/\/+$/, "")}/api`;
  }
  return (process.env.NEXT_PUBLIC_API_URL || "/api").replace(/\/+$/, "");
}

export async function apiFetch<T>(
  path: string,
  init?: RequestInit
): Promise<ApiCallResult<T>> {
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  const url = `${getApiBaseUrl()}${normalizedPath}`;
  const startedAt = performance.now();

  const res = await fetch(url, {
    ...init,
    cache: "no-store",
    headers: {
      Accept: "application/json",
      ...init?.headers,
    },
  });

  const latencyMs = Math.round(performance.now() - startedAt);
  const envelope = (await res.json()) as ApiEnvelope<T>;
  const requestId =
    res.headers.get("X-Request-ID") || envelope.meta?.request_id || null;

  if (!res.ok || !envelope.success) {
    const msg = envelope.error?.message || `Request failed with status ${res.status}`;
    throw new Error(msg);
  }

  return {
    envelope,
    status: res.status,
    latencyMs,
    requestId,
  };
}

export function fetchHealth() {
  return apiFetch<HealthData>("/v1/health");
}

export function fetchReady() {
  return apiFetch<ReadyData>("/v1/ready");
}

export function fetchSystemInfo() {
  return apiFetch<SystemInfoData>("/v1/info");
}
