<script lang="ts">
  export let event: CalendarEventModel;
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div {
    background-color: #cbe6ec;
    border-radius: $borderRadiusSmall;
    padding: $gapSmall;
  }
</style>

<div>
  {event.title}
</div>