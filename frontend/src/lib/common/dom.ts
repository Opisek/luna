export let generateUniqueElementId = (name: string[], randomness?: string | number | unknown) => (
  `${name.length === 0 ? "generic" : name.reverse().join("-")}.${randomness === undefined ? Math.floor(Math.random() * 100000000) : randomness}`
);
export let extendUniqueElementId = (name: string[], id?: string | unknown) => generateUniqueElementId(name, id);