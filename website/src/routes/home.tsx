import { createRoute } from "@tanstack/react-router";
import { motion } from "motion/react";
import { RootRoute } from "./root";
import { Hero } from "../components/Hero";
import { FeatureGrid } from "../components/FeatureGrid";
import { CatalogSection } from "../components/CatalogSection";
import { InstallSection } from "../components/InstallSection";
import { ApiOverview } from "../components/ApiOverview";
import { StarCard } from "../components/StarCard";
import { AiStack } from "../components/AiStack";

export const HomeRoute = createRoute({
  getParentRoute: () => RootRoute,
  path: "/",
  component: HomePage,
});
function HomePage() {
  return (
    <main>
      <Hero />
      <motion.section
        className="border-y border-line bg-[rgba(9,22,36,0.66)]"
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
      >
        <div className="mx-auto grid min-h-[246px] w-[calc(100%-48px)] max-w-[1160px] grid-cols-[1fr_1px_1fr] items-center gap-[55px] max-[680px]:block max-[680px]:py-16 max-[680px]:w-[calc(100%-34px)]">
          <div>
            <span className="mb-[15px] block font-mono text-[10px] text-mint">01</span>
            <h2 className="mb-3.5 max-w-[430px] text-[clamp(29px,3.5vw,44px)] leading-[1.1] tracking-[-0.065em]">
              A calm layer over native UI
            </h2>
            <p className="m-0 max-w-[500px] text-[13px] text-muted">
              Keep the speed and control of native rendering while giving your team a consistent
              design language.
            </p>
          </div>
          <div className="h-20 bg-line max-[680px]:hidden" />
          <div className="relative pl-2.5 max-[680px]:mt-[47px] max-[680px]:pl-[21px]">
            <span className="absolute top-[-29px] left-[-22px] font-serif text-[75px] leading-none text-mint">
              “
            </span>
            <p className="mb-3 text-base leading-[1.5] text-ink">
              Design tokens should make the right thing easy.
            </p>
            <small className="font-mono text-[9px] text-dim">— Breeze UI principles</small>
          </div>
        </div>
      </motion.section>
      <FeatureGrid />
      <AiStack />
      <CatalogSection />
      <ApiOverview />
      <InstallSection />
      <StarCard />
    </main>
  );
}
