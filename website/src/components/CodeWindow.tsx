import { Check, Circle } from "lucide-react";

export function CodeWindow() {
  return (
    <div className="overflow-hidden rounded-xl border border-[#234252] bg-[#091827] shadow-[0_28px_80px_rgba(0,0,0,0.38)]">
      <div className="flex h-[42px] items-center gap-2 border-b border-line px-4 font-mono text-[10px] text-[#688187]">
        <div className="flex gap-[5px] text-[#3c5260] [&>svg:first-child]:text-[#e59472] [&>svg:nth-child(2)]:text-[#eac37d] [&>svg:nth-child(3)]:text-[#6eca95]">
          <Circle size={8} fill="currentColor" />
          <Circle size={8} fill="currentColor" />
          <Circle size={8} fill="currentColor" />
        </div>
        <span>main.go</span>
        <span className="flex-1" />
      </div>
      <pre className="m-0 whitespace-pre-wrap p-6 font-mono text-[10px] leading-[2.02] text-[#bed0cf] max-[680px]:px-[14px] max-[680px]:py-[18px] max-[680px]:text-[8.5px] max-[680px]:leading-[1.95]">
        <span className="text-[#d3baff]">package</span> main{`\n\n`}
        <span className="text-[#d3baff]">import</span> ({`\n  `}
        <span className="text-[#9debc9]">"github.com/egoist/mygo/ui"</span>
        {`\n  `}breeze{" "}
        <span className="text-[#9debc9]">"github.com/yldm-tech/mygo-breeze-ui/ui/breeze"</span>
        {`\n`}){`\n\n`}
        <span className="text-[#d3baff]">func</span> <span className="text-[#ffc18f]">view</span>(c
        *ui.Context) {"{"}
        {`\n  `}breeze.<span className="text-[#ffc18f]">Card</span>(c,{" "}
        <span className="text-[#d3baff]">func</span>() {"{"}
        {`\n    `}breeze.<span className="text-[#ffc18f]">Text</span>(c,{" "}
        <span className="text-[#9debc9]">"Welcome home"</span>){`\n    `}
        <span className="text-[#d3baff]">if</span> breeze.
        <span className="text-[#ffc18f]">Button</span>(c,{" "}
        <span className="text-[#9debc9]">"Continue"</span>, breeze.ButtonPrimary).
        <span className="text-[#ffc18f]">Clicked</span>() {"{"}
        {`\n      `}
        <span className="text-[#668087]">// handle the action</span>
        {`\n    `}
        {"}"}
        {`\n  `}
        {"}"}){`\n`}
        {"}"}
      </pre>
      <div className="flex items-center gap-[7px] border-t border-line px-[17px] py-[11px] font-mono text-[9px] text-[#779095] [&>svg:first-child]:text-mint-strong">
        <Check size={13} />
        <span>compiled with Go</span>
        <span className="flex-1" />
        <span>native renderer</span>
      </div>
    </div>
  );
}
