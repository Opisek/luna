<style lang="scss">
    @import "../../styles/animations.scss";

    .wrapper {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        animation: spinner $animationSpeedSlow linear infinite;
    }

    @keyframes spinner {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }
</style>

<script lang="ts">
  import { LoaderCircle, Loader, RefreshCw } from "lucide-svelte";
</script>

<span class="wrapper">
    <LoaderCircle size={16}/>
    <!--<Loader/>-->
    <!--<RefreshCw/>-->
</span>