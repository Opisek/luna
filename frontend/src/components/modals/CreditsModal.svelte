<script lang="ts">
  import { CaseSensitive, Layers, Link, Package, Palette } from "lucide-svelte";
  import ButtonList from "../forms/ButtonList.svelte";
  import Modal from "./Modal.svelte";
  import { NoOp } from "$lib/client/placeholders";
  import type { Option } from "../../types/options";
  import IconButton from "../interactive/IconButton.svelte";
  import { Github } from "svelte-simples";
  import List from "../forms/List.svelte";
  import Vertical from "../layout/Vertical.svelte";

  interface Props {
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(),
  }: Props = $props();

  let showModalInternal: () => any = $state(NoOp);
  let hideModalInternal: () => any = $state(NoOp);;

  showModal = () => {
    showModalInternal();
  };

  hideModal = () => {
    hideModalInternal();
  };

  const categories: Option<string>[][] = [
    [
      //{ name: "Tech Stack", value: "languages", icon: Languages },
      { name: "Tech Stack", value: "languages", icon: Layers },
    ],
    [
      { name: "Frontend", value: "frontend", icon: Package },
      { name: "Themes", value: "themes", icon: Palette },
      { name: "Fonts", value: "fonts", icon: CaseSensitive },
    ],
    [
      { name: "Backend", value: "backend", icon: Package },
    ],
  ]

  let selectedCategory: string = $state("languages");
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.container {
    box-sizing: border-box;
    display: grid;
    grid-template-columns: auto 1fr;
    grid-template-rows: 1fr;
    gap: dimensions.$gapMiddle;
    min-width: 30vw;
    height: 60vh;
  }

  main {
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapMiddle;
    overflow-y: auto;
    overflow-x: hidden;
    padding-right: dimensions.$gapLarger;
    margin-right: -(dimensions.$gapLarger);
  }

  main > :global(*) {
    flex-shrink: 0;
  }

  .credit {
    padding: dimensions.$gapMiddle;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    border-radius: dimensions.$borderRadius;

    display: grid;
    gap: dimensions.$gapSmall;
    row-gap: 0;
    grid-template-columns: 1fr auto;
    grid-template-rows: auto auto;
    grid-template-areas: "name buttons" "details buttons";
    justify-content: center;
    align-items: center;
  }

  .credit > .name {
    grid-area: name;
    display: flex;
    flex-direction: row;
    justify-content: start;
    align-items: center;
    gap: dimensions.$gapTiny;
  }
  .credit > .details {
    grid-area: details;
    font-size: text.$fontSizeSmall;
  }
  .credit > .buttons {
    grid-area: buttons;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
  }
</style>

<Modal
  title={"Credits"}
  bind:showModal={showModalInternal}
  bind:hideModal={hideModalInternal}
>
  <div class="container">
    <Vertical height="auto" position="top">
      <ButtonList
        bind:value={selectedCategory}
        options={categories} 
      />
      <IconButton href="https://github.com/Opisek/luna">
        <Github/>
      </IconButton>
    </Vertical>
    <main tabindex="-1">
      {#if selectedCategory === "languages"}
        <List
          label="Tech Stack and Languages"
          items={[
            { name: "Bun", desc: "JavaScript Runtime", license: "MIT", url: "https://bun.com/", author: "Oven"},
            { name: "ECMAScript", desc: "Frontend Language", license: "Ecma Text Copyright Policy", url: "https://ecma-international.org/technical-committees/tc39/", author: "Ecma International"},
            { name: "Go", desc: "Backend Language", license: "Go License", url: "https://go.dev/", author: "Google"},
            { name: "PostgreSQL", desc: "Database Management System", license: "PostgreSQL License", url: "https://www.postgresql.org/", author: "The PostgreSQL Global Development Group"},
            { name: "Svelte", desc: "Frontend framework", license: "MIT", url: "https://svelte.dev", author: "Svelte contributors" },
            { name: "SvelteKit", desc: "Framework for Svelte", license: "MIT", url: "https://kit.svelte.dev", author: "Svelte contributors" },
            { name: "TypeScript", desc: "JavaScript Types", license: "Apache 2.0", url: "https://www.typescriptlang.org/", author: "Microsoft"},
          ]}
          id={item => item.url}
          template={creditTemplate}
        />
      {:else if selectedCategory === "frontend"}
        <List
          label="Frontend Libraries"
          items={[
            { name: "@types", desc: "Various TypeScript type definitions", license: "MIT", url: "https://github.com/DefinitelyTyped/DefinitelyTyped/", author: "DefinitelyTyped contributors"},
            { name: "Lucide (for Svelte)", desc: "Main icons library", license: "ISC", url: "https://lucide.dev/", author: "Lucide contributors"},
            { name: "Svelte Language Tools", desc: "Language tools for the Svelte framework", license: "MIT", url: "https://github.com/sveltejs/language-tools", author: "Svelte contributors"},
            { name: "Svelte Preprocess", desc: "Svelte preprocessor", license: "MIT", url: "https://github.com/sveltejs/svelte-preprocess", author: "Svelte contributors"},
            { name: "Svelte Simples", desc: "Secondary icons library", license: "MIT", url: "https://github.com/shinokada/svelte-simples", author: "Shinichi Okada"},
            { name: "UAParser.js", desc: "Parsing browser user agents", license: "MIT", url: "https://uaparser.dev/", author: "Faisal Salman"},
            { name: "Vite", desc: "Development web server", license: "MIT", url: "https://vite.dev/", author: "VoidZero Inc. & Vite contributors"},
            { name: "dotenv", desc: "Parsing of .env file", license: "BSD-2-Clause", url: "https://github.com/motdotla/dotenv", author: "motdotla"},
            { name: "iplocation", desc: "Determing login geolocation", license: "MIT", url: "https://github.com/Richienb/iplocation", author: "Richie Bendall"},
            { name: "node-sha1", desc: "SHA1 implementation", license: "BSD-3-Clause", url: "https://github.com/pvorb/node-sha1", author: "Paul Vorbach"},
            { name: "sha256", desc: "SHA256 implementation", license: "MIT", url: "https://github.com/cryptocoinjs/sha256", author: "JP Richardson"},
            { name: "svelte-adapter-bun", desc: "Svelte integration for bun web server", license: "MIT", url: "https://github.com/gornostay25/svelte-adapter-bun", author: "Volodymyr Palamar"},
            { name: "vite-plugin-svelte", desc: "Svelte integration for vite web server", license: "MIT", url: "https://github.com/sveltejs/vite-plugin-svelte", author: "Svelte contributors"},
          ]}
          id={item => item.url}
          template={creditTemplate}
        />
      {:else if selectedCategory === "backend"}
        <List
          label="Backend Libraries"
          items={[
            { name: "Cron", desc: "Cron job scheduler", license: "MIT", url: "https://github.com/robfig/cron", author: "Rob Figueiredo"},
            { name: "Gin", desc: "Web server framework", license: "MIT", url: "https://gin-gonic.com/", author: "Gin contributors"},
            { name: "GoDotEnv", desc: "Parsing .env files", license: "MIT", url: "https://github.com/joho/godotenv", author: "Brandon Keepers / John Barton"},
            { name: "Logrus", desc: "Structured logger", license: "MIT", url: "https://github.com/sirupsen/logrus", author: "Simon Eskildsen"},
            { name: "UUID (Google)", desc: "UUID generator (used for UUIDv4)", license: "BSD-3-Clause", url: "https://github.com/google/uuid", author: "Google"},
            { name: "UUID (The Go Commune)", desc: "UUID generator (used for UUIDv5)", license: "MIT", url: "https://github.com/gofrs/uuid", author: "The Go Commune"},
            { name: "crypto", desc: "Cryptographic library", license: "Go License", url: "https://pkg.go.dev/golang.org/x/crypto", author: "Google"},
            { name: "go-ical", desc: "Parsing iCal files", license: "MIT", url: "https://github.com/emersion/go-ical", author: "Simon Ser"},
            { name: "go-qrcode", desc: "QR code generator", license: "MIT", url: "https://github.com/skip2/go-qrcode", author: "skip2"},
            { name: "go-webdav", desc: "WebDAV protocol implementation", license: "MIT", url: "https://github.com/emersion/go-webdav", author: "Simon Ser"},
            { name: "golang-jwt", desc: "JSON Web Token parser and serializer", license: "MIT", url: "https://github.com/golang-jwt/jwt", author: "Michael Fridman"},
            { name: "pgx", desc: "PostgreSQL driver and toolkit", license: "MIT", url: "https://github.com/jackc/pgx", author: "Jack Christensen"},
            { name: "rrule-go", desc: "Date recurrence rules parser", license: "MIT", url: "https://github.com/teambition/rrule-go", author: "Teambition"},
            { name: "useragent", desc: "Parsing browser user agents", license: "MIT", url: "https://github.com/mileusna/useragent", author: "Misoš Mileusnić"},
          ]}
          id={item => item.url}
          template={creditTemplate}
        />
      {:else if selectedCategory === "fonts"}
        <List
          label="Pre-installed Fonts"
          items={[
            { name: "Atkinson Hyperlegible (Next/Mono)", desc: "Accessibility-focused typefaces", license: "SIL OPEN FONT LICENSE Version 1.1", url: "https://www.brailleinstitute.org/freefont/", author: "Braille Institute of America, Inc."},
          ]}
          id={item => item.url}
          template={creditTemplate}
        />
      {:else if selectedCategory === "themes"}
        <List
          label="Theme Color Definitions"
          items={[
            { name: "Catpuccin", desc: "", license: "MIT", url: "https://catppuccin.com/", author: "Catpuccin"},
            { name: "Dracula", desc: "", license: "MIT", url: "https://draculatheme.com/", author: "Zeno Rocha & Lucas de França"},
            { name: "Nord", desc: "", license: "MIT", url: "https://www.nordtheme.com/", author: "Arctic Ice Studio & Sven Greb"},
            { name: "Solarized", desc: "", license: "MIT", url: "https://ethanschoonover.com/solarized/", author: "Ethan Schoonover"},
          ]}
          id={item => item.url}
          template={creditTemplate}
        />
      {/if}
    </main>
  </div>
</Modal>

{#snippet creditTemplate(c: Credit)}
  <div class="credit">
    <span class="name">
      {c.name}
    </span>

    <span class="details">
      {#if c.desc === ""}
        {c.license}
      {:else}
        {c.desc} • {c.license}
      {/if}
    </span>

    <div class="buttons">
      <IconButton href={c.url}>
        {#if c.url.startsWith("https://github.com/")}
          <Github size={20}/>
        {:else}
          <Link size={20}/>
        {/if}
      </IconButton>
    </div>
  </div>
{/snippet}