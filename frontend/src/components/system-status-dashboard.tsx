"use client";

import { useCallback, useEffect, useState } from "react";
import {
  ActivityIcon,
  AlertCircleIcon,
  CheckCircle2Icon,
  CpuIcon,
  DatabaseIcon,
  GlobeIcon,
  RefreshCwIcon,
  ServerIcon,
  ShieldCheckIcon,
} from "lucide-react";
import { toast } from "sonner";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { fetchHealth, fetchReady, fetchSystemInfo } from "@/lib/api";
import type { HealthData, ReadyData, SystemInfoData } from "@/types/api";

interface TelemetryState {
  loading: boolean;
  error: string | null;
  health: HealthData | null;
  ready: ReadyData | null;
  info: SystemInfoData | null;
  latencyMs: number | null;
  requestId: string | null;
  checkedAt: string | null;
}

const initialState: TelemetryState = {
  loading: true,
  error: null,
  health: null,
  ready: null,
  info: null,
  latencyMs: null,
  requestId: null,
  checkedAt: null,
};

const apiEndpoints = [
  { method: "GET", path: "/nginx-health", target: "Nginx", desc: "Reverse proxy liveness probe" },
  { method: "GET", path: "/health", target: "Backend", desc: "Direct container liveness check" },
  { method: "GET", path: "/api/v1/health", target: "Backend", desc: "Versioned service health status" },
  { method: "GET", path: "/api/v1/ready", target: "Backend + DB", desc: "Readiness probe (HTTP server & PostgreSQL ping)" },
  { method: "GET", path: "/api/v1/info", target: "Backend + DB", desc: "Go 1.26 & PostgreSQL runtime metadata" },
];

function formatUptime(seconds: number): string {
  if (seconds < 60) return `${seconds}s`;
  const mins = Math.floor(seconds / 60);
  const remSec = seconds % 60;
  if (mins < 60) return `${mins}m ${remSec}s`;
  const hours = Math.floor(mins / 60);
  return `${hours}h ${mins % 60}m`;
}

export function SystemStatusDashboard() {
  const [state, setState] = useState<TelemetryState>(initialState);
  const [refreshing, setRefreshing] = useState(false);

  const fetchTelemetrySnapshot = useCallback(async () => {
    const [healthRes, readyRes, infoRes] = await Promise.all([
      fetchHealth(),
      fetchReady(),
      fetchSystemInfo(),
    ]);

    return {
      health: healthRes.envelope.data ?? null,
      ready: readyRes.envelope.data ?? null,
      info: infoRes.envelope.data ?? null,
      latencyMs: infoRes.latencyMs,
      requestId: infoRes.requestId,
      checkedAt: new Date().toLocaleTimeString(),
    };
  }, []);

  const handleRefresh = useCallback(async () => {
    setRefreshing(true);
    try {
      const snapshot = await fetchTelemetrySnapshot();
      setState({
        loading: false,
        error: null,
        ...snapshot,
      });
      toast.success("Backend telemetry refreshed", {
        description: `Responded in ${snapshot.latencyMs}ms (${snapshot.info?.go_version ?? "Go 1.26"})`,
      });
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Unable to reach backend API";
      setState((prev) => ({
        ...prev,
        loading: false,
        error: message,
        checkedAt: new Date().toLocaleTimeString(),
      }));
      toast.error("Backend check failed", { description: message });
    } finally {
      setRefreshing(false);
    }
  }, [fetchTelemetrySnapshot]);

  useEffect(() => {
    let cancelled = false;

    fetchTelemetrySnapshot()
      .then((snapshot) => {
        if (cancelled) return;
        setState({
          loading: false,
          error: null,
          ...snapshot,
        });
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        const message =
          err instanceof Error ? err.message : "Unable to reach backend API";
        setState((prev) => ({
          ...prev,
          loading: false,
          error: message,
          checkedAt: new Date().toLocaleTimeString(),
        }));
      });

    return () => {
      cancelled = true;
    };
  }, [fetchTelemetrySnapshot]);

  const isHealthy = !state.error && state.health?.status === "ok";
  const isDbConnected = !state.error && Boolean(state.info?.database?.connected);

  return (
    <div className="space-y-6">
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <div>
          <h2 className="text-xl font-semibold tracking-tight">
            System Telemetry
          </h2>
          <p className="text-sm text-muted-foreground">
            Live verification of Nginx → Go 1.26 Backend → PostgreSQL 18 connectivity.
          </p>
        </div>

        <Button
          variant="outline"
          onClick={() => void handleRefresh()}
          disabled={refreshing}
        >
          <RefreshCwIcon className={refreshing ? "animate-spin" : ""} />
          Refresh Status
        </Button>
      </div>

      {state.error && (
        <Alert variant="destructive">
          <AlertCircleIcon />
          <AlertTitle>Backend API Unreachable</AlertTitle>
          <AlertDescription>{state.error}</AlertDescription>
        </Alert>
      )}

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ServerIcon className="size-4 text-muted-foreground" />
              Backend API
            </CardTitle>
            <CardDescription>Golang 1.26 · net/http</CardDescription>
            <CardAction>
              {state.loading ? (
                <Skeleton className="h-5 w-16 rounded-full" />
              ) : (
                <Badge variant={isHealthy ? "default" : "destructive"}>
                  {isHealthy ? (
                    <>
                      <CheckCircle2Icon className="size-3" />
                      Healthy
                    </>
                  ) : (
                    "Offline"
                  )}
                </Badge>
              )}
            </CardAction>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            {state.loading ? (
              <div className="space-y-2">
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-4 w-1/2" />
              </div>
            ) : (
              <>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Runtime</span>
                  <span className="font-mono font-medium">
                    {state.info?.go_version ?? "—"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Uptime</span>
                  <span className="font-mono">
                    {state.health
                      ? formatUptime(state.health.uptime_seconds)
                      : "—"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Latency</span>
                  <span className="font-mono">
                    {state.latencyMs !== null ? `${state.latencyMs}ms` : "—"}
                  </span>
                </div>
              </>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <DatabaseIcon className="size-4 text-muted-foreground" />
              PostgreSQL
            </CardTitle>
            <CardDescription>pgx/v5 Pool · SQL Migrations</CardDescription>
            <CardAction>
              {state.loading ? (
                <Skeleton className="h-5 w-20 rounded-full" />
              ) : (
                <Badge variant={isDbConnected ? "default" : "destructive"}>
                  {isDbConnected ? (
                    <>
                      <CheckCircle2Icon className="size-3" />
                      Connected
                    </>
                  ) : (
                    "Disconnected"
                  )}
                </Badge>
              )}
            </CardAction>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            {state.loading ? (
              <div className="space-y-2">
                <Skeleton className="h-4 w-3/4" />
                <Skeleton className="h-4 w-1/2" />
              </div>
            ) : (
              <>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Engine</span>
                  <span className="font-mono font-medium">
                    {state.info?.database?.version ?? "PostgreSQL"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Probe</span>
                  <span className="font-mono">
                    {state.ready?.checks?.postgres ?? "—"}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Network</span>
                  <span className="font-mono">postgres:5432</span>
                </div>
              </>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <GlobeIcon className="size-4 text-muted-foreground" />
              Frontend UI
            </CardTitle>
            <CardDescription>Next.js 16 · Shadcn UI</CardDescription>
            <CardAction>
              <Badge variant="secondary">Active</Badge>
            </CardAction>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Router</span>
              <span className="font-mono">App Router</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Styling</span>
              <span className="font-mono">Tailwind v4</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Output</span>
              <span className="font-mono">Standalone</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ShieldCheckIcon className="size-4 text-muted-foreground" />
              Edge Proxy
            </CardTitle>
            <CardDescription>Nginx Alpine · Reverse Proxy</CardDescription>
            <CardAction>
              <Badge variant="outline">Port 80</Badge>
            </CardAction>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">API</span>
              <span className="font-mono">/api/* → :8080</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Web</span>
              <span className="font-mono">{"/* → :3000"}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Features</span>
              <span className="font-mono">Gzip · Limit</span>
            </div>
          </CardContent>
        </Card>
      </div>

      <Tabs defaultValue="runtime">
        <TabsList>
          <TabsTrigger value="runtime">
            <CpuIcon className="size-4" />
            Runtime Metadata
          </TabsTrigger>
          <TabsTrigger value="endpoints">
            <ActivityIcon className="size-4" />
            System Endpoints
          </TabsTrigger>
        </TabsList>

        <TabsContent value="runtime">
          <Card>
            <CardHeader>
              <CardTitle>Backend & Database Diagnostics</CardTitle>
              <CardDescription>
                Live payload returned by <code className="font-mono">GET /api/v1/info</code> and{" "}
                <code className="font-mono">GET /api/v1/ready</code>.
              </CardDescription>
            </CardHeader>
            <CardContent className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
              {[
                { label: "Service", value: state.info?.service },
                { label: "Version", value: state.info?.version },
                { label: "Environment", value: state.info?.environment },
                {
                  label: "Platform",
                  value: state.info ? `${state.info.os}/${state.info.arch} (${state.info.num_cpu} vCPU)` : undefined,
                },
                {
                  label: "Database",
                  value: state.info?.database?.version ?? state.info?.database?.engine,
                },
                {
                  label: "Readiness Checks",
                  value: state.ready
                    ? Object.entries(state.ready.checks)
                        .map(([k, v]) => `${k}:${v}`)
                        .join(" · ")
                    : undefined,
                },
              ].map((item) => (
                <div
                  key={item.label}
                  className="rounded-lg border bg-muted/30 p-3"
                >
                  <p className="text-xs text-muted-foreground">{item.label}</p>
                  {state.loading ? (
                    <Skeleton className="mt-1 h-5 w-24" />
                  ) : (
                    <p className="mt-1 font-mono text-sm font-medium">
                      {item.value ?? "—"}
                    </p>
                  )}
                </div>
              ))}
            </CardContent>
            <CardFooter className="justify-between text-xs text-muted-foreground">
              <span>
                Request ID:{" "}
                <code className="font-mono">{state.requestId ?? "—"}</code>
              </span>
              <span>Last checked: {state.checkedAt ?? "—"}</span>
            </CardFooter>
          </Card>
        </TabsContent>

        <TabsContent value="endpoints">
          <Card>
            <CardHeader>
              <CardTitle>Exposed Base Routes</CardTitle>
              <CardDescription>
                All routes are served through the single Nginx entrypoint.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              {apiEndpoints.map((ep, idx) => (
                <div key={ep.path}>
                  <div className="flex flex-wrap items-center justify-between gap-2 py-1">
                    <div className="flex items-center gap-2">
                      <Badge variant="secondary" className="font-mono">
                        {ep.method}
                      </Badge>
                      <code className="font-mono text-sm font-medium">
                        {ep.path}
                      </code>
                    </div>
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <span>{ep.desc}</span>
                      <Badge variant="outline">{ep.target}</Badge>
                    </div>
                  </div>
                  {idx < apiEndpoints.length - 1 && <Separator className="mt-2" />}
                </div>
              ))}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
