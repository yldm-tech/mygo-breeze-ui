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
    <section className="hero container">
      <motion.div
        className="hero-copy"
        initial={{ opacity: 0, y: 18 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.65 }}
      >
        <p className="eyebrow">
          <span className="pulse" />
          {t.eyebrow}
        </p>
        <h1>{t.title}</h1>
        <p className="hero-intro">{t.intro}</p>
        <div className="hero-actions">
          <a className="button button-primary" href="#install">
            {t.start}
            <ArrowRight size={16} />
          </a>
          <a className="button button-ghost" href="#catalog">
            {t.explore}
          </a>
        </div>
        <div className="stats">
          <div>
            <strong>{components}</strong>
            <span>{t.components}</span>
          </div>
          <div>
            <strong>{utilities}</strong>
            <span>{t.utilities}</span>
          </div>
          <div>
            <strong>0</strong>
            <span>{t.dependencies}</span>
          </div>
        </div>
      </motion.div>
      <motion.div
        className="hero-visual"
        initial={{ opacity: 0, x: 22 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ duration: 0.75, delay: 0.12 }}
      >
        <CodeWindow />
        <div className="float-card float-card-top">
          <Sparkles size={16} />
          <span>typed · composable · native</span>
        </div>
        <div className="float-card float-card-bottom">
          <span className="float-check">✓</span>
          <span>no DOM required</span>
        </div>
      </motion.div>
    </section>
  );
}

export function HeroLinks() {
  return (
    <div className="hero-links">
      <Link to="/api">
        API reference <ArrowUpRight size={14} />
      </Link>
      <a href="https://github.com/yldm-tech/mygo-breeze-ui" target="_blank" rel="noreferrer">
        GitHub <ArrowUpRight size={14} />
      </a>
    </div>
  );
}
