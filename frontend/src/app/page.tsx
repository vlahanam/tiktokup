import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { SystemStatusDashboard } from "@/components/system-status-dashboard";

export default function HomePage() {
  return (
    <div className="space-y-8">
      <section className="space-y-3">
        <div className="flex flex-wrap items-center gap-2">
          <Badge>Go 1.26</Badge>
          <Badge variant="secondary">PostgreSQL 18</Badge>
          <Badge variant="secondary">Next.js 16 + TypeScript</Badge>
          <Badge variant="secondary">Shadcn UI</Badge>
          <Badge variant="outline">Nginx Reverse Proxy</Badge>
          <Badge variant="outline">Docker Compose</Badge>
        </div>

        <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">
          TikTokUp Base Platform
        </h1>
        <p className="max-w-2xl text-muted-foreground">
          Production-ready Dockerized monorepo foundation featuring a clean Go
          1.26 HTTP API, PostgreSQL 18 database, Next.js App Router + Shadcn UI
          frontend, and an Nginx reverse proxy entrypoint.
        </p>
      </section>

      <Separator />

      <SystemStatusDashboard />
    </div>
  );
}
