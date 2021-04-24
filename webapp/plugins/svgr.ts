import fs from "fs";
import type { Plugin } from "vite";
import type * as E from "esbuild";

export default function svgrPlugin(): Plugin {
  return {
    name: "vite:svgr",
    async transform(code, id) {
      if (id.endsWith(".svg")) {
        const svgr = require("@svgr/core").default;
        const esbuild = require("esbuild") as typeof E;

        const svg = await fs.promises.readFile(id, "utf8");

        const componentCode = (
          await svgr(svg, { ref: true }, { componentName: "ReactComponent" })
        ).replace(
          "export default ForwardRef;",
          "export { ForwardRef as ReactComponent };"
        );

        const res = await esbuild.transform(componentCode + "\n" + code, {
          loader: "jsx",
        });

        return {
          code: res.code,
        };
      }
    },
  };
}
