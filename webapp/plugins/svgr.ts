import fs from "fs";
import type { Plugin } from "vite";
import type * as E from "esbuild";

export default function svgrPlugin(): Plugin {
  return {
    name: "vite:svgr",
    async transform(code, id) {
      if (id.endsWith(".svg")) {
        const { transform } = require("@svgr/core");
        const esbuild = require("esbuild") as typeof E;

        const svg = await fs.promises.readFile(id, "utf8");

        const componentCode = (
          await transform(
            svg,
            {
              ref: true,
              plugins: ["@svgr/plugin-svgo", "@svgr/plugin-jsx"],
              svgoConfig: {
                plugins: [
                  {
                    name: "preset-default",
                    params: {
                      overrides: {
                        removeViewBox: false,
                      },
                    },
                  },
                ],
              },
            },
            { componentName: "ReactComponent" }
          )
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
