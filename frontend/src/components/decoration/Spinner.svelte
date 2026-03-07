<script lang="ts">
  import { LoaderCircle } from "lucide-svelte";
</script>

<style lang="scss">
    @use "../../styles/animations.scss";

    .wrapper {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        animation: spin animations.$animationSpeedSlow linear infinite;
    }

    .wrapper::before {
        content: "a"; // This forces the spinner to take the same amount of space as text would.
        visibility: hidden;
    }
</style>

<span class="wrapper">
    <LoaderCircle size={16}/>
    <!--<Loader/>-->
    <!--<RefreshCw/>-->
</span>