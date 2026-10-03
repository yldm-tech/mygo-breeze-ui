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
    <section className="features container">
      {features.map(({ icon: Icon, color, title, text, code }, index) => (
        <motion.article
          className="feature"
          key={title}
          initial={{ opacity: 0, y: 18 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, amount: 0.4 }}
          transition={{ delay: index * 0.08 }}
        >
          <div className={`feature-icon ${color}`}>
            <Icon size={22} />
          </div>
          <h3>{t[title]}</h3>
          <p>{t[text]}</p>
          <code>{code}</code>
        </motion.article>
      ))}
    </section>
  );
}
