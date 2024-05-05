<script lang="ts">
  export let calendar: CalendarModel;
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
    width: 1em;
    height: 1em;
    border-radius: 50%;
  }

  span {
    flex-grow: 1;
  }
</style>

<div class="entry">
  <div class="color" style="background-color:{calendar.color}"></div>
  <span>
    {calendar.name}
  </span>
  <input type="checkbox">
</div>