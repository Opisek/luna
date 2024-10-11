<script lang="ts">
  export let value: boolean;
  export let description: string;
  export let name: string;

  export let editable: boolean = true;

  export let onChange: (value: boolean) => any = () => {};
</script>

<style lang="scss">
  div {
    display: flex;
    align-items: start;
    flex-direction: row;
    flex-wrap: nowrap;
  }
</style>

<div>
  <input
    bind:checked={value}
    on:change={() => onChange(value)}
    type="checkbox"
    disabled={!editable}
    name={name}
  />
  <label for={name}>
    {description}
  </label>
</div>