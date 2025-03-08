<script lang="ts">
  import { LoaderCircle } from "lucide-svelte";
</script>

<style lang="scss">
    @import "../../styles/animations.scss";

    .wrapper {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        animation: spin $animationSpeedSlow linear infinite;
    }
</style>

<span class="wrapper">
    <LoaderCircle size={16}/>
    <!--<Loader/>-->
    <!--<RefreshCw/>-->
</span>