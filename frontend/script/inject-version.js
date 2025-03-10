import {readFileSync, writeFileSync} from "fs";

const version = process.argv[2]

function replaceInFile(path, replaceFunc) {
    let contents = readFileSync(path, "utf-8");
    contents = replaceFunc(contents);
    writeFileSync(path, contents);
}

replaceInFile("../package.json", (contents) => {
    const pack = JSON.parse(contents);
    pack.version = version;
    return JSON.stringify(pack, null, 2);
})

replaceInFile("../src/lib/common/version.ts", (contents) => {
    return contents.replace(/\"[^\"]+\"/, `"${version}"`);
})