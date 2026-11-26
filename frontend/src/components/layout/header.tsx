import Link from "next/link";
import { LayersIcon } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { ThemeToggle } from "@/components/theme-toggle";

const appName = process.env.NEXT_PUBLIC_APP_NAME || "TikTokUp";

export function Header() {
  return (
    <header className="sticky top-0 z-40 border-b bg-background/80 backdrop-blur-md">
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4 sm:px-6">
        <Link href="/" className="flex items-center gap-2.5 font-semibold">
          <span className="flex size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <LayersIcon className="size-4" />
          </span>
          <span className="tracking-tight">{appName}</span>
          <Badge variant="secondary" className="hidden sm:inline-flex">
            Base v0.1
          </Badge>
        </Link>

        <div className="flex items-center gap-2">
          <Badge variant="outline" className="hidden md:inline-flex">
            Docker Stack
          </Badge>
          <ThemeToggle />
        </div>
      </div>
    </header>
  );
}
