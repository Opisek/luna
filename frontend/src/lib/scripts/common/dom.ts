export let generateUniqueElementId = (name: string[], randomness?: string | number | unknown) => (
  `${name.length === 0 ? "generic" : name.join("_")}-${randomness === undefined ? Math.floor(Math.random() * 100000000) : randomness}`
);
export let extendUniqueElementId = (name: string[], id?: string | unknown) => generateUniqueElementId(name, id);