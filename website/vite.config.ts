import { copyFile, mkdir } from "node:fs/promises";
import { defineConfig, lazyPlugins } from "vite-plus";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

const referenceComponents = "Accordion,AccordionItem,ActionRow,AlertDialog,Autocomplete,Avatar,BackButton,Badge,Breadcrumbs,Button,Calendar,CalendarView,Card,Checkbox,CheckboxGroup,Collapsible,ColorPicker,ColorWell,ComboRow,Combobox,DateInput,DatePicker,DecimalInput,Divider,EditableText,EntryRow,ExpanderRow,Field,Fieldset,FindBar,Form,ForwardButton,GridView,HeaderBar,IconButton,InfoBar,Link,List,MenuButton,Meter,Modal,NumberInput,Outline,OutlineTable,PasswordField,Popover,PopoverMenuButton,PrimaryButton,Progress,ProgressRing,Radio,RadioGroup,RangeSlider,Rating,Scroll,ScrollBoth,ScrollHorizontal,SearchField,Segmented,Select,Sidebar,SidebarItem,SidebarSection,Slider,SpinRow,Spinner,Split,SplitButton,SplitPercent,SplitVertical,SplitVerticalPercent,StatusPage,StepSlider,Stepper,Switch,SwitchRow,TabView,Table,Tabs,Text,TextArea,TextAreaBuffer,TextInput,TextInputBuffer,TimeInput,Toasts,Toggle,ToggleGroup,TokenField,Toolbar,Tooltip,Tree,TreeItem,ViewSwitcher".split(",");

const githubPagesFallback = () => ({
  name: "github-pages-fallback",
  closeBundle: async () => {
    await copyFile("dist/index.html", "dist/404.html");
    await mkdir("dist/docs", { recursive: true });
    await copyFile("dist/index.html", "dist/docs/index.html");
    await mkdir("dist/docs/components", { recursive: true });
    await copyFile("dist/index.html", "dist/docs/components/index.html");
    await Promise.all(referenceComponents.map(async (component) => {
      const directory = `dist/docs/components/${component}`;
      await mkdir(directory, { recursive: true });
      await copyFile("dist/index.html", `${directory}/index.html`);
    }));
    await mkdir("dist/cli-mcp", { recursive: true });
    await copyFile("dist/index.html", "dist/cli-mcp/index.html");
  },
});

export default defineConfig({
  base: "/",
  plugins: lazyPlugins(() => [react(), tailwindcss(), githubPagesFallback()]),
  server: { port: 4173 },
});
