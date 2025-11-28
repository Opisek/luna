import { error, json } from "@sveltejs/kit";
import type { RequestEvent } from "./$types";
import * as fs from "fs";
import path from "path";
import { apiProxy } from "$lib/server/api.server";

const validResources = [ "themes", "fonts" ];
const validResourceFileNameRegex = new RegExp(/[a-zA-Z0-9-]+\.[a-zA-Z0-9]+/);

// This endpoint determines what files are inside /static/resource/ directory.
// This is for telling the UI which fonts or themes are available.
export const GET = async ({ params, request, url, getClientAddress }: RequestEvent) => {
  // Make sure the user is logged in
  const response = await apiProxy(request, getClientAddress(), "sessions/permissions", { method: "GET", body: null }, false);
  if (!response.ok) return response;

  // Verify that the requested resource is valid
  const requestedResource = params.resource;
  if (!validResources.includes(requestedResource)) {
    return error(400, "Invalid resource type");
  }

  // List the files
  const resourceDir = process.env.DEVELOPMENT == "true" ? `static` : `build/client`;
  const resourcePath = `./${resourceDir}/${requestedResource}`;
  const resources = recursivelyFindFiles(resourceDir, resourcePath);

  return json(resources);
};

// This endpoint adds files to the appropriate /static/resource/ directory.
// It can only be used by administrators.
export const PUT = async ({ params, request, url, getClientAddress }: RequestEvent) => {
  // Make sure the user is logged in and an admin
  const response = await apiProxy(request, getClientAddress(), "sessions/permissions", { method: "GET", body: null }, false);
  if (!response.ok) return response;
  const data = (await response.json()) as { user_id: string, is_admin: boolean, permissions: string[] };
  if (!data.is_admin || !data.permissions.includes('administrative') || !data.permissions.includes('manage_resources')) return error(401, "Unauthorized");

  // Verify that the requested resource is valid
  const requestedResource = params.resource;
  if (!validResources.includes(requestedResource)) {
    return error(400, "Invalid resource type");
  }

  // Extract the file
  const formData = await request.formData();
  if (!formData.has("file")) return error(400, "Missing resource file");
  const file = formData.get("file");
  if (file == null || !(file instanceof File)) return error(400, "Missing or malformed resource file");

  // Verify valid file name
  if (!validResourceFileNameRegex.test(file.name)) error(400, "File name contains illegal characters");

  // Verify correct file type
  const fileExtension = file.name.split(".")[1];
  if (!(() => {
    switch (requestedResource) {
      case "themes":
        if (fileExtension != "css") return false;
        break;
      case "fonts":
        if (fileExtension != "ttf") return false;
        break;
    }
    return true;
  })()) return error(400, "Invalid file extension");

  // Save the file
  const resourceDir = process.env.DEVELOPMENT == "true" ? `static` : `build/client`;
  let resourcePath = `./${resourceDir}/${requestedResource}`;
  if (requestedResource == "themes") {
    const body = await file.text()
    if (body.includes("html[data-theme=\"dark\"]")) resourcePath += "/dark";
    else if (body.includes("html[data-theme=\"light\"]")) resourcePath += "/light";
    else return error(400, "Theme file does not specify whether it is a light or dark theme");
  }

  fs.writeFileSync(safelyJoinPath(resourcePath, file.name), await file.text());

  return json({ "status": "ok" });
};

// This endpoint deletes files from the appropriate /static/resource/ directory.
// It can only be used by administrators.
export const DELETE = async ({ params, request, url, getClientAddress }: RequestEvent) => {
  // Make sure the user is logged in and an admin
  const response = await apiProxy(request, getClientAddress(), "sessions/permissions", { method: "GET", body: null }, false);
  if (!response.ok) return response;
  const data = (await response.json()) as { user_id: string, is_admin: boolean, permissions: string[] };
  if (!data.is_admin || !data.permissions.includes('administrative') || !data.permissions.includes('manage_resources')) return error(401, "Unauthorized");

  // Verify that the requested resource is valid
  const resourceParts = params.resource.split("/");
  if (resourceParts.length < 2) return error(400, "Invalid resource path");
  const requestedResource = resourceParts.shift();
  if (!validResources.includes(requestedResource as string)) {
    return error(400, "Invalid resource type");
  }
  const remainingResourcePath = resourceParts.join("/");

  // Check if the file exists
  const resourceDir = process.env.DEVELOPMENT == "true" ? `static` : `build/client`;
  const resourceRootPath = `./${resourceDir}/${requestedResource}`;
  let resourcePath = safelyJoinPath(resourceRootPath, remainingResourcePath);

  // Add the appropriate extension
  switch (requestedResource) {
    case "themes":
      resourcePath += ".css";
      break;
    case "fonts":
      resourcePath += ".ttf";
      break;
  }

  if (!fs.existsSync(resourcePath)) return error(404, "Resource not found");
  const stat = fs.statSync(resourcePath);
  if (!stat.isFile()) return error(404, "Resource not found");

  fs.unlinkSync(resourcePath);

  return json({ "status": "ok" });
};

// https://stackoverflow.com/questions/37862886/safe-path-resolve-without-escaping-from-parent
function safelyJoinPath(base: string, target: string): string {
  const targetPath = "." + path.normalize("/" + target);
  return path.resolve(base, targetPath);
}

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