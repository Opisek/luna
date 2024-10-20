<script lang="ts">
  export let href: string;
  export let onClick: () => any = () => {};
</script>

<style lang="scss">
  @import "../../styles/colors.scss";

  a {
    text-decoration: none;
    color: $foregroundLink;
    width: max-content;
  }
</style>

<a href={href} on:click={onClick()}><slot/></a>