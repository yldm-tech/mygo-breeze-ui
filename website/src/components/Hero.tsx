import { motion } from "motion/react";
import { ArrowRight, ArrowUpRight, Sparkles } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { getCatalog } from "../lib/catalog";
import { useUiStore } from "../store";
import { useCopy } from "../copy";
import { CodeWindow } from "./CodeWindow";

export function Hero() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  const { data } = useQuery({ queryKey: ["catalog"], queryFn: getCatalog });
  const components = data?.components.length ?? 11;
  const utilities = data?.utilities.length ?? 33;
  return (
    <section className="mx-auto grid min-h-[650px] w-[calc(100%-48px)] max-w-[1160px] grid-cols-[1.02fr_0.98fr] items-center gap-[66px] py-[68px] pb-[98px] max-[940px]:grid-cols-1 max-[940px]:gap-[45px] max-[680px]:min-h-0 max-[680px]:w-[calc(100%-34px)] max-[680px]:py-[55px] max-[680px]:pb-[73px]">
      <motion.div
        className="relative z-[1]"
        initial={{ opacity: 0, y: 18 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.65 }}
      >
        <p className="mb-5 flex items-center gap-[9px] font-mono text-[10px] uppercase tracking-[0.09em] text-mint">
          <span className="size-[7px] rounded-full bg-mint-strong shadow-[0_0_0_5px_rgba(98,230,177,0.1)]" />
          {t.eyebrow}
        </p>
        <h1 className="mb-[26px] max-w-[680px] text-[clamp(44px,5.3vw,74px)] font-bold leading-[1.04] tracking-[-0.078em]">
          {t.title}
        </h1>
        <p className="mb-[31px] max-w-[565px] text-[15px] leading-[1.85] text-muted max-[680px]:text-[13px]">
          {t.intro}
        </p>
        <div className="flex items-center gap-2.5">
          <a
            className="inline-flex items-center justify-center gap-2 rounded-[7px] bg-mint px-[17px] py-[11px] text-xs font-extrabold text-bg transition hover:-translate-y-0.5 hover:bg-[#c0f9df]"
            href="#install"
          >
            {t.start}
            <ArrowRight size={16} />
          </a>
          <a
            className="inline-flex items-center justify-center gap-2 rounded-[7px] border border-line bg-transparent px-[17px] py-[10px] text-xs font-extrabold text-ink transition hover:-translate-y-0.5 hover:border-muted"
            href="#catalog"
          >
            {t.explore}
          </a>
        </div>
        <div className="mt-[63px] flex gap-[31px] max-[680px]:mt-[43px] max-[680px]:gap-[19px]">
          <div className="flex flex-col gap-0.5">
            <strong className="text-[21px] tracking-[-0.05em] max-[680px]:text-[18px]">
              {components}
            </strong>
            <span className="font-mono text-[9px] uppercase tracking-[0.08em] text-dim">
              {t.components}
            </span>
          </div>
          <div className="flex flex-col gap-0.5">
            <strong className="text-[21px] tracking-[-0.05em] max-[680px]:text-[18px]">
              {utilities}
            </strong>
            <span className="font-mono text-[9px] uppercase tracking-[0.08em] text-dim">
              {t.utilities}
            </span>
          </div>
          <div className="flex flex-col gap-0.5">
            <strong className="text-[21px] tracking-[-0.05em] max-[680px]:text-[18px]">0</strong>
            <span className="font-mono text-[9px] uppercase tracking-[0.08em] text-dim">
              {t.dependencies}
            </span>
          </div>
        </div>
      </motion.div>
      <motion.div
        className="relative"
        initial={{ opacity: 0, x: 22 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.75, delay: 0.12 }}
      >
        <CodeWindow />
        <div className="absolute top-[34px] right-0 flex items-center gap-2 rounded-[7px] border border-[#345264] bg-[rgba(17,35,55,0.9)] px-[11px] py-[9px] font-mono text-[9px] text-muted shadow-[0_14px_35px_rgba(0,0,0,0.25)] backdrop-blur-xl max-[940px]:right-[-5px] max-[680px]:hidden">
          <Sparkles size={16} />
          <span>typed · composable · native</span>
        </div>
        <div className="absolute bottom-12 left-[-31px] flex items-center gap-2 rounded-[7px] border border-[#345264] bg-[rgba(17,35,55,0.9)] px-[11px] py-[9px] font-mono text-[9px] text-muted shadow-[0_14px_35px_rgba(0,0,0,0.25)] backdrop-blur-xl max-[940px]:left-[-5px] max-[680px]:hidden">
          <span className="grid size-4 place-items-center rounded-full bg-mint text-[10px] text-bg">
            ✓
          </span>
          <span>no DOM required</span>
        </div>
      </motion.div>
    </section>
  );
}

export function HeroLinks() {
  return (
    <div className="flex gap-4 font-mono text-[10px] text-mint">
      <Link className="inline-flex items-center gap-1.5" to="/api">
        API reference <ArrowUpRight size={14} />
      </Link>
      <a
        className="inline-flex items-center gap-1.5"
        href="https://github.com/yldm-tech/mygo-breeze-ui"
        target="_blank"
        rel="noreferrer"
      >
        GitHub <ArrowUpRight size={14} />
      </a>
    </div>
  );
}
