import { Separator } from "@/components/ui/separator";

export function Footer() {
  return (
    <footer className="mt-auto">
      <Separator />
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-2 px-4 py-5 text-xs text-muted-foreground sm:flex-row sm:px-6">
        <p>TikTokUp Base Platform — Containerized Full-Stack Architecture</p>
        <p className="font-mono">Next.js 16 · TypeScript · Shadcn UI · Go 1.26 · PostgreSQL 18 · Nginx</p>
      </div>
    </footer>
  );
}
