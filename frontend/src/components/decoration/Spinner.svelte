<script lang="ts">
  import { LoaderCircle } from "lucide-svelte";
</script>

<style lang="scss">
    @use "../../styles/animations.scss";

    .spinner {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        animation: spin animations.$animationSpeedSlow linear infinite;
    }
</style>

<span class="spinner">
    <LoaderCircle size={16}/>
    <!--<Loader/>-->
    <!--<RefreshCw/>-->
</span>