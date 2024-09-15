import {readFileSync, writeFileSync} from "fs";

const version = process.argv[2]
const pack = JSON.parse(readFileSync("../package.json", "utf-8"))
pack.version = version
writeFileSync("../package.json", JSON.stringify(pack, null, 2))