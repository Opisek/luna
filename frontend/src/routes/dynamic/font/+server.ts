export const GET = ({ url }) => {
  const file = url.searchParams.get("file") || "";
  const name = file.split("-").map(x => x.charAt(0).toUpperCase() + x.slice(1)).join(" ");
  const purpose = url.searchParams.get("purpose");

  const stylesheet = `
    @font-face {
      font-family: "${name}"; 
      src: url("/fonts/${file}.ttf");
    }

    html {
      --${purpose}: ${name};
    }
  `;

  return new Response(stylesheet, {
    headers: [
      ["Content-Type", "text/css"]
    ]
  });
}