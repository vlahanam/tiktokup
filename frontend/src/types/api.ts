export interface ApiMeta {
  request_id?: string;
  timestamp: string;
}

export interface ApiErrorInfo {
  code: string;
  message: string;
}

export interface ApiEnvelope<T> {
  success: boolean;
  data?: T;
  error?: ApiErrorInfo;
  meta: ApiMeta;
}

export interface HealthData {
  status: string;
  service: string;
  environment: string;
  uptime_seconds: number;
}

export interface ReadyData {
  status: string;
  checks: Record<string, string>;
}

export interface DatabaseInfo {
  engine: string;
  connected: boolean;
  version?: string;
}

export interface SystemInfoData {
  service: string;
  version: string;
  environment: string;
  go_version: string;
  os: string;
  arch: string;
  num_cpu: number;
  started_at: string;
  uptime_seconds: number;
  database?: DatabaseInfo;
}
