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
        className="principles"
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        viewport={{ once: true }}
      >
        <div className="container principles-grid">
          <div>
            <span className="section-number">01</span>
            <h2>A calm layer over native UI</h2>
            <p>
              Keep the speed and control of native rendering while giving your team a consistent
              design language.
            </p>
          </div>
          <div className="principle-line" />
          <div className="principle-quote">
            <span>“</span>
            <p>Design tokens should make the right thing easy.</p>
            <small>— Breeze UI principles</small>
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
