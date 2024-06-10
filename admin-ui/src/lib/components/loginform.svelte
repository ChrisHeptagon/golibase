<script lang="ts">
  export let data;
  export let type: "Login" | "Register" = "Register";
</script>

<section>
  <h1>{type}</h1>
  <form>
    {#each Object.keys(data) as key}
      <label for={key}>{key}</label>
      <input type={data[key].type} id={key} autocomplete="on" />
    {/each}
    <button>{type}</button>
  </form>
</section>

<style lang="scss">
  @font-face {
    font-family: "Lora";
    src: url("$lib/fonts/lora.woff2") format("woff2");
  }
section {
    height: 500px;
    width: 500px;
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    background-color: gray;
    font-family: "Lora", serif;
}
h1 {
    font-size: 2em;
    color: white;
    text-align: center;
    
}
form {
    padding: 10px;
    display: flex;
    flex-direction: column;
    justify-content: center;
}
label {
    font-size: 1.5em;
    color: white;
    padding: 5px;
}
input {
    font-size: 1.5em;
    padding: 5px;
    margin: 5px;
    border-radius: 5px;
}
button {
    font-size: 1.5em;
    padding: 5px;
    background-color: white;
    color: black;
    border-radius: 5px;
    margin: 5px;
}
</style>
