import { Check, Circle } from "lucide-react";

export function CodeWindow() {
  return (
    <div className="code-window">
      <div className="window-bar">
        <div className="window-dots">
          <Circle size={8} fill="currentColor" />
          <Circle size={8} fill="currentColor" />
          <Circle size={8} fill="currentColor" />
        </div>
        <span>main.go</span>
        <span className="window-spacer" />
      </div>
      <pre>
        <span className="code-keyword">package</span> main{`\n\n`}
        <span className="code-keyword">import</span> ({`\n  `}
        <span className="code-string">"github.com/egoist/mygo/ui"</span>
        {`\n  `}breeze{" "}
        <span className="code-string">"github.com/yldm-tech/mygo-breeze-ui/ui/breeze"</span>
        {`\n`}){`\n\n`}
        <span className="code-keyword">func</span> <span className="code-function">view</span>(c
        *ui.Context) {"{"}
        {`\n  `}breeze.<span className="code-function">Card</span>(c,{" "}
        <span className="code-keyword">func</span>() {"{"}
        {`\n    `}breeze.<span className="code-function">Text</span>(c,{" "}
        <span className="code-string">"Welcome home"</span>){`\n    `}
        <span className="code-keyword">if</span> breeze.
        <span className="code-function">Button</span>(c,{" "}
        <span className="code-string">"Continue"</span>, breeze.ButtonPrimary).
        <span className="code-function">Clicked</span>() {"{"}
        {`\n      `}
        <span className="code-comment">// handle the action</span>
        {`\n    `}
        {"}"}
        {`\n  `}
        {"}"}){`\n`}
        {"}"}
      </pre>
      <div className="code-status">
        <Check size={13} />
        <span>compiled with Go</span>
        <span className="window-spacer" />
        <span>native renderer</span>
      </div>
    </div>
  );
}
