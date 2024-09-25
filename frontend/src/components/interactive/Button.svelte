<script lang="ts">
  export let onClick: () => void;

  // TODO: could not figure out enums for this, try again later
  export let color: string;
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  button {
    // unset props
    outline: none;
    border: 0;
    margin: 0;

    cursor: pointer;
    padding: $gapSmall;
    border-radius: $borderRadius;

    min-width: 5em;
  }


  // TODO: might change to an scss loop later
  button.success {
    background-color: $backgroundSuccess;
    color: $foregroundSuccess;
  }
  button.failure {
    background-color: $backgroundFailure;
    color: $foregroundFailure;
  }
  button.accent {
    background-color: $backgroundAccent;
    color: $foregroundAccent;
  }
</style>

<button
  on:click={onClick}
  class:success={color == "success"}
  class:failure={color == "failure"}
  class:accent={color == "accent"}
>
  <slot></slot>
</button>