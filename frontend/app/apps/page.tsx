"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { listWorkflows, type SavedWorkflowSummary } from "@/lib/api";

export default function AppsListPage() {
  const [items, setItems] = useState<SavedWorkflowSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    listWorkflows()
      .then((data) => {
        if (!cancelled) {
          setItems(data);
          setLoading(false);
        }
      })
      .catch((e) => {
        if (!cancelled) {
          setError(e instanceof Error ? e.message : String(e));
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <main className="flex-1 overflow-y-auto p-8">
      <div className="max-w-3xl mx-auto space-y-6">
        <div>
          <h1 className="text-xl font-semibold">Saved mini-apps</h1>
          <p className="text-sm text-zinc-500 mt-1">
            Workflows you've saved. Open one to re-run with new inputs.
          </p>
        </div>

        {loading && <p className="text-sm text-zinc-500">Loading…</p>}
        {error && (
          <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
        )}
        {!loading && !error && items.length === 0 && (
          <p className="text-sm text-zinc-500 italic">
            No saved mini-apps yet.{" "}
            <Link href="/" className="underline">
              Go compose one
            </Link>
            .
          </p>
        )}

        <ul className="space-y-2">
          {items.map((a) => (
            <li
              key={a.id}
              className="rounded-md border border-black/10 dark:border-white/15 p-3 hover:border-black/40 dark:hover:border-white/40 transition-colors"
            >
              <Link href={`/apps/${a.id}`} className="block">
                <div className="font-medium text-sm">{a.name}</div>
                {a.description && (
                  <div className="text-xs text-zinc-600 dark:text-zinc-400 mt-0.5">
                    {a.description}
                  </div>
                )}
                <div className="text-[11px] text-zinc-500 mt-1">
                  {a.step_count} step{a.step_count !== 1 ? "s" : ""} ·{" "}
                  {a.created_at
                    ? new Date(a.created_at).toLocaleString()
                    : "—"}{" "}
                  · <span className="font-mono">{a.id.slice(0, 8)}</span>
                </div>
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </main>
  );
}
