import { motion } from "motion/react";
import { Layers3, Sparkles, Terminal } from "lucide-react";
import { useUiStore } from "../store";
import { useCopy } from "../copy";

const features = [
  {
    icon: Sparkles,
    color: "mint",
    title: "utilityTitle",
    text: "utilityText",
    code: "breeze.Row(c, breeze.Gap(4), breeze.P(6))",
  },
  {
    icon: Layers3,
    color: "orange",
    title: "componentTitle",
    text: "componentText",
    code: 'breeze.Dialog(c, &open, "Hello", body, footer)',
  },
  {
    icon: Terminal,
    color: "violet",
    title: "toolingTitle",
    text: "toolingText",
    code: "breeze-ui add Button -o button.go",
  },
] as const;

export function FeatureGrid() {
  const { locale } = useUiStore();
  const t = useCopy(locale);
  return (
    <section className="mx-auto grid w-[calc(100%-48px)] max-w-[1160px] grid-cols-3 gap-[22px] py-24 pb-[108px] max-[680px]:w-[calc(100%-34px)] max-[680px]:grid-cols-1 max-[680px]:gap-[46px] max-[680px]:py-[68px] max-[680px]:pb-20">
      {features.map(({ icon: Icon, color, title, text, code }, index) => (
        <motion.article
          className="border-t border-line pt-[22px]"
          key={title}
          initial={{ opacity: 0, y: 18 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, amount: 0.4 }}
          transition={{ delay: index * 0.08 }}
        >
          <div
            className={`mb-[25px] grid size-[43px] place-items-center rounded-[11px] ${color === "mint" ? "bg-[rgba(157,245,208,0.08)] text-mint" : color === "orange" ? "bg-[rgba(255,174,118,0.08)] text-orange" : "bg-[rgba(184,166,255,0.08)] text-violet"}`}
          >
            <Icon size={22} />
          </div>
          <h3 className="mb-2 text-base tracking-[-0.025em]">{t[title]}</h3>
          <p className="mb-[18px] min-h-[67px] text-xs leading-[1.75] text-muted max-[680px]:min-h-0">
            {t[text]}
          </p>
          <code className="block max-w-full overflow-hidden text-ellipsis whitespace-nowrap font-mono text-[9px] text-dim">
            {code}
          </code>
        </motion.article>
      ))}
    </section>
  );
}
