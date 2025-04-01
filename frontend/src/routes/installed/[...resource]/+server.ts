import { error, json } from "@sveltejs/kit";
import type { RequestEvent } from "./$types";
import * as fs from "fs";
import path from "path";

const validResources = [ "themes", "fonts" ];

// This endpoint determines what files are inside /static/resource/ directory
// This is for telling the UI which fonts or themes are available
export const GET = async ({ params, request, url }: RequestEvent) => {
  const requestedResource = params.resource;

  // Verify that the requested resource is valid
  if (!validResources.includes(requestedResource)) {
    return error(400, "Invalid resource type");
  }

  // Open the /static/resource directory
  const resourceDir = process.env.DEVELOPMENT == "true" ? `static` : `build/client`;
  const resourcePath = `./${resourceDir}/${requestedResource}`;
  const resources = recursivelyFindFiles(resourceDir, resourcePath);

  return json(resources);
};

function recursivelyFindFiles(baseDir: string, dir: string) {
  const found = {} as any;

  const files = fs.readdirSync(dir);
  for (const file of files) {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);

    if (stat.isDirectory()) {
      found[file] = recursivelyFindFiles(baseDir, filePath);
    } else {
      found[file.split(".")[0]] = filePath.split(".")[0].substring(baseDir.length);
    }
  }

  return found;
}