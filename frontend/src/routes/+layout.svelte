<style lang="scss">
  @import "../styles/dimensions.scss";

  :global(*) {
    box-sizing: border-box;
  }

  :global(body) {
    margin: 0;
    padding: 0;
    height: 100vh;
    font-family: 'Karla', serif;
    font-size: $fontSize;
  }
</style>

<slot/>