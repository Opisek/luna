import { register, init } from "@sveltia/i18n";
import { parse } from "yaml";

register("en-US", () => import("../../lang/en-US.yaml?raw").then(m => parse(m.default)));

init({ fallbackLocale: "en-US" });