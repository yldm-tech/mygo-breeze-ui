import { ArrowRight, Braces, Database, HeartPulse } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { apiEndpoints } from "../lib/api";
import { ApiPlayground } from "./ApiPlayground";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

const icons = [Database, Braces, HeartPulse];
export function ApiOverview() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <section className="px-0 py-[100px] pb-[108px] max-[680px]:py-[75px] max-[680px]:pb-20">
      <div className="mx-auto w-[calc(100%-48px)] max-w-[1160px] max-[680px]:w-[calc(100%-34px)]">
        <div className="mb-9 flex items-end justify-between gap-7 max-[680px]:block">
          <div>
            <p className="mb-5 font-mono text-[10px] uppercase tracking-[0.09em] text-mint">
              {t.apiKicker}
            </p>
            <h2 className="mb-[9px] text-[clamp(29px,3.5vw,44px)] leading-[1.1] tracking-[-0.065em]">
              {t.apiTitle}
            </h2>
            <p className="m-0 text-[13px] text-muted">{t.apiText}</p>
          </div>
          <Link
            className="inline-flex items-center gap-1.5 text-xs font-extrabold text-mint hover:text-ink max-[680px]:mt-5"
            to="/api"
          >
            {t.apiLink}
            <ArrowRight size={15} />
          </Link>
        </div>
        <div className="grid grid-cols-3 gap-2.5 max-[940px]:grid-cols-2 max-[680px]:grid-cols-1">
          {apiEndpoints.map((endpoint, index) => {
            const Icon = icons[index] ?? Database;
            return (
              <article
                className="flex gap-3.5 rounded-lg border border-line bg-panel p-5"
                key={endpoint.path}
              >
                <div className="grid size-[34px] shrink-0 place-items-center rounded-lg bg-[rgba(157,245,208,0.08)] text-mint">
                  <Icon size={17} />
                </div>
                <div>
                  <div className="font-mono text-[10px] text-dim">
                    <span className="mr-1.5 text-mint">{endpoint.method}</span>
                    {endpoint.path}
                  </div>
                  <h3 className="mb-2 mt-2 text-base tracking-[-0.025em]">{endpoint.label}</h3>
                  <p className="m-0 text-[11px] leading-[1.6] text-muted">{endpoint.description}</p>
                </div>
              </article>
            );
          })}
        </div>
        <ApiPlayground />
      </div>
    </section>
  );
}
