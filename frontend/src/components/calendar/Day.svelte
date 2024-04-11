<script lang="ts">
  export let day: number;
  export let dayOfWeek: number;
  export let isCurrentMonth: boolean;
</script>

<style lang="scss">
  div.otherMonth {
    opacity: .5;
  }
  div.sunday {
    color: red;
  }
</style>

<div class:otherMonth={!isCurrentMonth} class:sunday={dayOfWeek === 0}>
  {day}
</div>