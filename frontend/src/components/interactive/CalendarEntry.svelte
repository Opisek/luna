<script lang="ts">
  import VisibilityToggle from "./VisibilityToggle.svelte";

  export let calendar: CalendarModel;

  let visibleTest = true;
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div.entry {
    display: flex;
    flex-direction: row;
    gap: $gapSmall;
    padding: $gapSmall;
    width: 100%;
    align-items: center;
  }

  div.color {
    width: 0.75em;
    height: 0.75em;
    border-radius: 50%;
  }

  span {
    flex-grow: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>

<div class="entry">
  <div class="color" style="background-color:{calendar.color}"></div>
  <span>
    {calendar.name}
  </span>
  <VisibilityToggle visible={visibleTest}/>
</div>